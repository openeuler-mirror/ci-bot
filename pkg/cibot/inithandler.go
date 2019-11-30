package cibot

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

type InitHandler struct {
	Config      config.Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

type Projects struct {
	Community    Community    `yaml:"community"`
	Repositories []Repository `yaml:"repositories"`
}

type Community struct {
	Name              *string  `yaml:"name"`
	ProtectedBranches []string `yaml:"protected_branches"`
	Managers          []string `yaml:"managers"`
	Developers        []string `yaml:"developers"`
	Viewers           []string `yaml:"viewers"`
	Reporters         []string `yaml:"reporters"`
}

type Repository struct {
	Name              *string  `yaml:"name"`
	Description       *string  `yaml:"description"`
	ProtectedBranches []string `yaml:"protected_branches"`
	Type              *string  `yaml:"type"`
	Managers          []string `yaml:"managers"`
	Developers        []string `yaml:"developers"`
	Viewers           []string `yaml:"viewers"`
	Reporters         []string `yaml:"reporters"`
}

var (
	PrivilegeManager   = "manager"
	PrivilegeDeveloper = "developer"
	PrivilegeViewer    = "viewer"
	PrivilegeReporter  = "reporter"

	PermissionAdmin = "admin"
	PermissionPush  = "push"
	PermissionPull  = "pull"

	BranchProtected = "protected"
	// not supported yet
	BranchReadonly = "readonly"
)

// Serve
func (handler *InitHandler) Serve() {
	// init waiting sha
	err := handler.initWaitingSha()
	if err != nil {
		glog.Errorf("unable to initWaitingSha: %v", err)
		return
	}
	// watch database
	handler.watch()
}

// initWaitingSha init waiting sha
func (handler *InitHandler) initWaitingSha() error {
	if len(handler.Config.WatchProjectFiles) == 0 {
		return nil
	}

	for _, wf := range handler.Config.WatchProjectFiles {
		// get params
		watchOwner := wf.WatchProjectFileOwner
		watchRepo := wf.WatchprojectFileRepo
		watchPath := wf.WatchprojectFilePath
		watchRef := wf.WatchProjectFileRef

		// invoke api to get file contents
		localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
		localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
		localVarOptionals.Ref = optional.NewString(watchRef)

		// get contents
		contents, _, err := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
			handler.Context, watchOwner, watchRepo, watchPath, localVarOptionals)
		if err != nil {
			glog.Errorf("unable to get repository content: %v", err)
			return err
		}
		// Check project file
		var lenProjectFiles int
		err = database.DBConnection.Model(&database.ProjectFiles{}).
			Where("owner = ? and repo = ? and path = ? and ref = ?", watchOwner, watchRepo, watchPath, watchRef).
			Count(&lenProjectFiles).Error
		if err != nil {
			glog.Errorf("unable to get project files: %v", err)
			return err
		}
		if lenProjectFiles > 0 {
			glog.Infof("project file is exist: %s", contents.Sha)
			// Check sha in database
			updatepf := database.ProjectFiles{}
			err = database.DBConnection.
				Where("owner = ? and repo = ? and path = ? and ref = ?", watchOwner, watchRepo, watchPath, watchRef).
				First(&updatepf).Error
			if err != nil {
				glog.Errorf("unable to get project files: %v", err)
				return err
			}
			// write sha in waitingsha
			updatepf.WaitingSha = contents.Sha
			err = database.DBConnection.Save(&updatepf).Error
			if err != nil {
				glog.Errorf("unable to get project files: %v", err)
				return err
			}

		} else {
			glog.Infof("project file is non-exist: %s", contents.Sha)
			// add project file
			addpf := database.ProjectFiles{
				Owner:      watchOwner,
				Repo:       watchRepo,
				Path:       watchPath,
				Ref:        watchRef,
				WaitingSha: contents.Sha,
			}

			// create project file
			err = database.DBConnection.Create(&addpf).Error
			if err != nil {
				glog.Errorf("unable to create project files: %v", err)
				return err
			}
			glog.Infof("add project file successfully: %s", contents.Sha)
		}
	}
	return nil
}

// watch database
func (handler *InitHandler) watch() {
	if len(handler.Config.WatchProjectFiles) == 0 {
		return
	}

	for {
		watchDuration := handler.Config.WatchProjectFileDuration
		for _, wf := range handler.Config.WatchProjectFiles {
			// get params
			watchOwner := wf.WatchProjectFileOwner
			watchRepo := wf.WatchprojectFileRepo
			watchPath := wf.WatchprojectFilePath
			watchRef := wf.WatchProjectFileRef

			glog.Infof("begin to serve. watchOwner: %s watchRepo: %s watchPath: %s watchRef: %s watchDuration: %d",
				watchOwner, watchRepo, watchPath, watchRef, watchDuration)

			// get project file
			pf := database.ProjectFiles{}
			err := database.DBConnection.
				Where("owner = ? and repo = ? and path = ? and ref = ?", watchOwner, watchRepo, watchPath, watchRef).
				First(&pf).Error
			if err != nil {
				glog.Errorf("unable to get project files: %v", err)
			} else {
				glog.Infof("init handler current sha: %v target sha: %v waiting sha: %v",
					pf.CurrentSha, pf.TargetSha, pf.WaitingSha)
				if pf.TargetSha != "" {
					// skip when there is executing target sha
					glog.Infof("there is executing target sha: %v", pf.TargetSha)
				} else {
					if pf.WaitingSha != "" && pf.CurrentSha != pf.WaitingSha {
						// waiting -> target
						pf.TargetSha = pf.WaitingSha
						err = database.DBConnection.Save(&pf).Error
						if err != nil {
							glog.Errorf("unable to save project files: %v", err)
						} else {
							// define update pf
							updatepf := &database.ProjectFiles{}
							updatepf.ID = pf.ID

							// get file content from target sha
							glog.Infof("get target sha blob: %v", pf.TargetSha)
							localVarOptionals := &gitee.GetV5ReposOwnerRepoGitBlobsShaOpts{}
							localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
							blob, _, err := handler.GiteeClient.GitDataApi.GetV5ReposOwnerRepoGitBlobsSha(
								handler.Context, watchOwner, watchRepo, pf.TargetSha, localVarOptionals)
							if err != nil {
								glog.Errorf("unable to get blob: %v", err)
							} else {
								// base64 decode
								glog.Infof("decode target sha blob: %v", pf.TargetSha)
								decodeBytes, err := base64.StdEncoding.DecodeString(blob.Content)
								if err != nil {
									glog.Errorf("decode content with error: %v", err)
								} else {
									// unmarshal owners file
									glog.Infof("unmarshal target sha blob: %v", pf.TargetSha)
									var ps Projects
									err = yaml.Unmarshal(decodeBytes, &ps)
									if err != nil {
										glog.Errorf("failed to unmarshal projects: %v", err)
									} else {
										glog.Infof("get blob result: %v", ps)
										result := true
										for i := 0; i < len(ps.Repositories); i++ {
											// get repositories length
											lenRepositories, err := handler.getRepositoriesLength(*ps.Community.Name, *ps.Repositories[i].Name)
											if err != nil {
												glog.Errorf("failed to get repositories length: %v", err)
												result = false
												continue
											}
											if lenRepositories > 0 {
												glog.Infof("repository: %s is exist. no action.", *ps.Repositories[i].Name)
											} else {
												// add repository
												err = handler.addRepositories(*ps.Community.Name, *ps.Repositories[i].Name,
													*ps.Repositories[i].Description, *ps.Repositories[i].Type)
												if err != nil {
													glog.Errorf("failed to add repositories: %v", err)
													result = false
													continue
												}
											}
											// add members
											err = handler.handleMembers(ps.Community, ps.Repositories[i])
											if err != nil {
												glog.Errorf("failed to add members: %v", err)
												result = false
											}
											// handle branches
											err = handler.handleBranches(ps.Community, ps.Repositories[i])
											if err != nil {
												glog.Errorf("failed to handle branches: %v", err)
												result = false
											}
											// handle repository type
											err = handler.handleRepositoryTypes(ps.Community, ps.Repositories[i])
											if err != nil {
												glog.Errorf("failed to handle repository types: %v", err)
												result = false
											}
										}
										glog.Infof("running result: %v", result)
										if result {
											err = database.DBConnection.Model(updatepf).Update("CurrentSha", pf.TargetSha).Error
											if err != nil {
												glog.Errorf("unable to update current sha: %v", err)
											}
										}
									}
								}
							}

							// at last update target sha
							err = database.DBConnection.Model(updatepf).Update("TargetSha", "").Error
							if err != nil {
								glog.Errorf("unable to update target sha: %v", err)
							}
							glog.Info("update sha successfully")
						}
					} else {
						glog.Infof("no waiting sha: %v", pf.WaitingSha)
					}
				}
			}
		}

		// watch duration
		glog.Info("end to serve")
		time.Sleep(time.Duration(watchDuration) * time.Second)
	}
}

// GetRepositoriesLength get repositories length
func (handler *InitHandler) getRepositoriesLength(owner string, repo string) (int, error) {
	// Check repositories file
	var lenRepositories int
	err := database.DBConnection.Model(&database.Repositories{}).
		Where("owner = ? and repo = ?", owner, repo).
		Count(&lenRepositories).Error
	if err != nil {
		glog.Errorf("unable to get repositories files: %v", err)
	}
	return lenRepositories, err
}

// addRepositories add repository
func (handler *InitHandler) addRepositories(owner, repo, description, t string) error {
	// add repository in gitee
	err := handler.addRepositoriesinGitee(owner, repo, description, t)
	if err != nil {
		glog.Errorf("failed to add repositories: %v", err)
		return err
	}

	// add repository in database
	err = handler.addRepositoriesinDB(owner, repo, description, t)
	if err != nil {
		glog.Errorf("failed to add repositories: %v", err)
		return err
	}
	return nil
}

// addRepositoriesinDB add repository in database
func (handler *InitHandler) addRepositoriesinDB(owner, repo, description, t string) error {
	// add repository
	addrepo := database.Repositories{
		Owner:       owner,
		Repo:        repo,
		Description: description,
		Type:        t,
	}

	// create repository
	err := database.DBConnection.Create(&addrepo).Error
	if err != nil {
		glog.Errorf("unable to create repository: %v", err)
		return err
	}
	return nil
}

// addRepositoriesinGitee add repository in giteee
func (handler *InitHandler) addRepositoriesinGitee(owner, repo, description, t string) error {
	// build create repository param
	repobody := gitee.RepositoryPostParam{}
	repobody.AccessToken = handler.Config.GiteeToken
	repobody.Name = repo
	repobody.Description = description
	repobody.HasIssues = true
	repobody.HasWiki = true
	// set `auto_init` as true to initialize `master` branch with README after repo creation
	repobody.AutoInit = true
	if t == "private" {
		repobody.Private = true
	} else {
		repobody.Private = false
	}

	// invoke query repository
	glog.Infof("begin to query repository: %s", repo)
	localVarOptionals := &gitee.GetV5ReposOwnerRepoOpts{}
	localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
	_, response, _ := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepo(handler.Context, owner, repo, localVarOptionals)
	if response.StatusCode == 404 {
		glog.Infof("repository is not exist: %s", repo)
	} else {
		glog.Infof("repository is already exist: %s", repo)
		return nil
	}

	// invoke create repository
	glog.Infof("begin to create repository: %s", repo)
	_, _, err := handler.GiteeClient.RepositoriesApi.PostV5OrgsOrgRepos(handler.Context, owner, repobody)
	if err != nil {
		glog.Errorf("fail to create repository: %v", err)
		return err
	}
	glog.Infof("end to create repository: %s", repo)
	return nil
}

// isUsingRepositoryMember check if using repository member not community member
func (handler *InitHandler) isUsingRepositoryMember(r Repository) bool {
	return len(r.Managers) > 0 || len(r.Developers) > 0 || len(r.Viewers) > 0 || len(r.Reporters) > 0
}

// handleMembers handle members
func (handler *InitHandler) handleMembers(c Community, r Repository) error {
	// if the members is defined in the repositories, it means that
	// all the members defined in the community will not be inherited by repositories.
	members := make(map[string]map[string]string)
	if handler.isUsingRepositoryMember(r) {
		// using repositories members
		glog.Infof("using repository members: %s", *r.Name)
		members = handler.getMembersMap(r.Managers, r.Developers, r.Viewers, r.Reporters)

	} else {
		// using community members
		glog.Infof("using community members: %s", *r.Name)
		members = handler.getMembersMap(c.Managers, c.Developers, c.Viewers, c.Reporters)
	}

	// get members from database
	var ps []database.Privileges
	err := database.DBConnection.Model(&database.Privileges{}).
		Where("owner = ? and repo = ?", c.Name, r.Name).Find(&ps).Error
	if err != nil {
		glog.Errorf("unable to get members: %v", err)
		return err
	}
	membersinDB := handler.getMembersMapByDB(ps)

	/* reporters
	err = handler.removeReporters(c, r, members[PrivilegeReporter], membersinDB[PrivilegeReporter])
	if err != nil {
		glog.Errorf("unable to remove reporters: %v", err)
	}*/
	// viewers
	err = handler.removeViewers(c, r, members[PrivilegeViewer], membersinDB[PrivilegeViewer])
	if err != nil {
		glog.Errorf("unable to remove viewers: %v", err)
	}
	// developers
	err = handler.removeDevelopers(c, r, members[PrivilegeDeveloper], membersinDB[PrivilegeDeveloper])
	if err != nil {
		glog.Errorf("unable to remove developers: %v", err)
	}
	// managers
	err = handler.removeManagers(c, r, members[PrivilegeManager], membersinDB[PrivilegeManager])
	if err != nil {
		glog.Errorf("unable to remove managers: %v", err)
	}

	/* currently reporters are not supported by gitee api
	err = handler.addReporters(c, r, members[PrivilegeReporter], membersinDB[PrivilegeReporter])
	if err != nil {
		glog.Errorf("unable to add reporters: %v", err)
	}*/
	// viewers
	err = handler.addViewers(c, r, members[PrivilegeViewer], membersinDB[PrivilegeViewer])
	if err != nil {
		glog.Errorf("unable to add viewers: %v", err)
	}
	// developers
	err = handler.addDevelopers(c, r, members[PrivilegeDeveloper], membersinDB[PrivilegeDeveloper])
	if err != nil {
		glog.Errorf("unable to add developers: %v", err)
	}
	// managers
	err = handler.addManagers(c, r, members[PrivilegeManager], membersinDB[PrivilegeManager])
	if err != nil {
		glog.Errorf("unable to add managers: %v", err)
	}

	return nil
}

// getMembersMap get members map
func (handler *InitHandler) getMembersMap(managers, developers, viewers, reporters []string) map[string]map[string]string {
	mapManagers := make(map[string]string)
	mapDevelopers := make(map[string]string)
	mapViewers := make(map[string]string)
	mapReporters := make(map[string]string)
	if len(managers) > 0 {
		for _, m := range managers {
			mapManagers[m] = m
		}
	}
	if len(developers) > 0 {
		for _, d := range developers {
			// skip when developer is already in managers
			_, okinManagers := mapManagers[d]
			if !okinManagers {
				mapDevelopers[d] = d
			}
		}
	}
	if len(viewers) > 0 {
		for _, v := range viewers {
			// skip when viewer is already in managers or developers
			_, okinManagers := mapManagers[v]
			_, okinDevelopers := mapDevelopers[v]
			if !okinManagers && !okinDevelopers {
				mapViewers[v] = v
			}
		}
	}
	if len(reporters) > 0 {
		for _, rt := range reporters {
			// skip when reporter is already in managers or developers or viewer
			_, okinManagers := mapManagers[rt]
			_, okinDevelopers := mapDevelopers[rt]
			_, okinViewers := mapViewers[rt]
			if !okinManagers && !okinDevelopers && !okinViewers {
				mapReporters[rt] = rt
			}
		}
	}

	// all members map
	members := make(map[string]map[string]string)
	members[PrivilegeManager] = mapManagers
	members[PrivilegeDeveloper] = mapDevelopers
	members[PrivilegeViewer] = mapViewers
	members[PrivilegeReporter] = mapReporters
	return members
}

// getMembersMapByDB get members map from database
func (handler *InitHandler) getMembersMapByDB(ps []database.Privileges) map[string]map[string]string {
	members := make(map[string]map[string]string)
	mapManagers := make(map[string]string)
	mapDevelopers := make(map[string]string)
	mapViewers := make(map[string]string)
	mapReporters := make(map[string]string)
	// all members map
	members[PrivilegeManager] = mapManagers
	members[PrivilegeDeveloper] = mapDevelopers
	members[PrivilegeViewer] = mapViewers
	members[PrivilegeReporter] = mapReporters
	if len(ps) > 0 {
		for _, p := range ps {
			members[p.Type][p.User] = strconv.Itoa(int(p.ID))
		}
	}

	return members
}

// addManagers add managers
func (handler *InitHandler) addManagers(c Community, r Repository, mapManagers, mapManagersInDB map[string]string) error {
	// managers added
	listOfAddManagers := make([]string, 0)
	for m := range mapManagers {
		_, okinManagers := mapManagersInDB[m]
		if !okinManagers {
			listOfAddManagers = append(listOfAddManagers, m)
		}
	}
	glog.Infof("list of add managers: %v", listOfAddManagers)
	if len(listOfAddManagers) > 0 {
		// build create project member param
		memberbody := gitee.ProjectMemberPutParam{}
		memberbody.AccessToken = handler.Config.GiteeToken
		memberbody.Permission = PermissionAdmin

		glog.Infof("begin to create manager for: %s", *r.Name)
		for j := 0; j < len(listOfAddManagers); j++ {
			_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfAddManagers[j], memberbody)
			if err != nil {
				glog.Errorf("fail to create manager: %v", err)
				continue
			}
			// create privilege
			ps := database.Privileges{
				Owner: *c.Name,
				Repo:  *r.Name,
				User:  listOfAddManagers[j],
				Type:  PrivilegeManager,
			}
			err = database.DBConnection.Create(&ps).Error
			if err != nil {
				glog.Errorf("fail to create manager in database: %v", err)
			}
		}
		glog.Infof("end to create manager for: %s", *r.Name)
	}
	return nil
}

// removeManagers remove managers
func (handler *InitHandler) removeManagers(c Community, r Repository, mapManagers, mapManagersInDB map[string]string) error {
	// managers removed
	listOfRemoveManagers := make([]string, 0)
	for m := range mapManagersInDB {
		_, okinManagers := mapManagers[m]
		if !okinManagers {
			listOfRemoveManagers = append(listOfRemoveManagers, m)
		}
	}
	glog.Infof("list of removed managers: %v", listOfRemoveManagers)
	if len(listOfRemoveManagers) > 0 {
		// build remove project member param
		memberbody := &gitee.DeleteV5ReposOwnerRepoCollaboratorsUsernameOpts{}
		memberbody.AccessToken = optional.NewString(handler.Config.GiteeToken)

		glog.Infof("begin to remove managers for: %s", *r.Name)
		for j := 0; j < len(listOfRemoveManagers); j++ {
			_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfRemoveManagers[j], memberbody)
			if err != nil {
				glog.Errorf("fail to remove managers: %v", err)
				continue
			}
			// delete privilege
			id, _ := strconv.Atoi(mapManagersInDB[listOfRemoveManagers[j]])
			ps := database.Privileges{}
			ps.ID = uint(id)
			err = database.DBConnection.Delete(&ps).Error
			if err != nil {
				glog.Errorf("fail to remove manager in database: %v", err)
			}
		}
		glog.Infof("end to remove managers for: %s", *r.Name)
	}
	return nil
}

// addDevelopers add developers
func (handler *InitHandler) addDevelopers(c Community, r Repository, mapDevelopers, mapDevelopersInDB map[string]string) error {
	// developers added
	listOfAddDevelopers := make([]string, 0)
	for d := range mapDevelopers {
		_, okinDevelopers := mapDevelopersInDB[d]
		if !okinDevelopers {
			listOfAddDevelopers = append(listOfAddDevelopers, d)
		}
	}
	glog.Infof("list of add developers: %v", listOfAddDevelopers)
	if len(listOfAddDevelopers) > 0 {
		// build create project member param
		memberbody := gitee.ProjectMemberPutParam{}
		memberbody.AccessToken = handler.Config.GiteeToken
		memberbody.Permission = PermissionPush

		glog.Infof("begin to create developers for: %s", *r.Name)
		for j := 0; j < len(listOfAddDevelopers); j++ {
			_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfAddDevelopers[j], memberbody)
			if err != nil {
				glog.Errorf("fail to create developers: %v", err)
				continue
			}
			// create privilege
			ps := database.Privileges{
				Owner: *c.Name,
				Repo:  *r.Name,
				User:  listOfAddDevelopers[j],
				Type:  PrivilegeDeveloper,
			}
			err = database.DBConnection.Create(&ps).Error
			if err != nil {
				glog.Errorf("fail to create developers in database: %v", err)
			}
		}
		glog.Infof("end to create developers for: %s", *r.Name)
	}
	return nil
}

// removeDevelopers remove developers
func (handler *InitHandler) removeDevelopers(c Community, r Repository, mapDevelopers, mapDevelopersInDB map[string]string) error {
	// developers removed
	listOfRemoveDevelopers := make([]string, 0)
	for d := range mapDevelopersInDB {
		_, okinDevelopers := mapDevelopers[d]
		if !okinDevelopers {
			listOfRemoveDevelopers = append(listOfRemoveDevelopers, d)
		}
	}
	glog.Infof("list of removed developers: %v", listOfRemoveDevelopers)
	if len(listOfRemoveDevelopers) > 0 {
		// build remove project member param
		memberbody := &gitee.DeleteV5ReposOwnerRepoCollaboratorsUsernameOpts{}
		memberbody.AccessToken = optional.NewString(handler.Config.GiteeToken)

		glog.Infof("begin to remove developers for: %s", *r.Name)
		for j := 0; j < len(listOfRemoveDevelopers); j++ {
			_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfRemoveDevelopers[j], memberbody)
			if err != nil {
				glog.Errorf("fail to remove developers: %v", err)
				continue
			}
			// delete privilege
			id, _ := strconv.Atoi(mapDevelopersInDB[listOfRemoveDevelopers[j]])
			ps := database.Privileges{}
			ps.ID = uint(id)
			err = database.DBConnection.Delete(&ps).Error
			if err != nil {
				glog.Errorf("fail to remove developers in database: %v", err)
			}
		}
		glog.Infof("end to remove developers for: %s", *r.Name)
	}
	return nil
}

// addViewers add viewers
func (handler *InitHandler) addViewers(c Community, r Repository, mapViewers, mapViewersInDB map[string]string) error {
	// viewers added
	listOfAddViewers := make([]string, 0)
	for v := range mapViewers {
		_, okinViewers := mapViewersInDB[v]
		if !okinViewers {
			listOfAddViewers = append(listOfAddViewers, v)
		}
	}
	glog.Infof("list of add viewers: %v", listOfAddViewers)
	if len(listOfAddViewers) > 0 {
		// build create project member param
		memberbody := gitee.ProjectMemberPutParam{}
		memberbody.AccessToken = handler.Config.GiteeToken
		memberbody.Permission = PermissionPull

		glog.Infof("begin to create viewers for: %s", *r.Name)
		for j := 0; j < len(listOfAddViewers); j++ {
			_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfAddViewers[j], memberbody)
			if err != nil {
				glog.Errorf("fail to create viewers: %v", err)
				continue
			}
			// create privilege
			ps := database.Privileges{
				Owner: *c.Name,
				Repo:  *r.Name,
				User:  listOfAddViewers[j],
				Type:  PrivilegeViewer,
			}
			err = database.DBConnection.Create(&ps).Error
			if err != nil {
				glog.Errorf("fail to create viewers in database: %v", err)
			}
		}
		glog.Infof("end to create viewers for: %s", *r.Name)
	}
	return nil
}

// removeViewers remove viewers
func (handler *InitHandler) removeViewers(c Community, r Repository, mapViewers, mapViewersInDB map[string]string) error {
	// viewers removed
	listOfRemoveViewers := make([]string, 0)
	for v := range mapViewersInDB {
		_, okinViewers := mapViewers[v]
		if !okinViewers {
			listOfRemoveViewers = append(listOfRemoveViewers, v)
		}
	}
	glog.Infof("list of removed viewers: %v", listOfRemoveViewers)
	if len(listOfRemoveViewers) > 0 {
		// build remove project member param
		memberbody := &gitee.DeleteV5ReposOwnerRepoCollaboratorsUsernameOpts{}
		memberbody.AccessToken = optional.NewString(handler.Config.GiteeToken)

		glog.Infof("begin to remove viewers for: %s", *r.Name)
		for j := 0; j < len(listOfRemoveViewers); j++ {
			_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfRemoveViewers[j], memberbody)
			if err != nil {
				glog.Errorf("fail to remove viewers: %v", err)
				continue
			}
			// delete privilege
			id, _ := strconv.Atoi(mapViewersInDB[listOfRemoveViewers[j]])
			ps := database.Privileges{}
			ps.ID = uint(id)
			err = database.DBConnection.Delete(&ps).Error
			if err != nil {
				glog.Errorf("fail to remove viewers in database: %v", err)
			}
		}
		glog.Infof("end to remove viewers for: %s", *r.Name)
	}
	return nil
}

// addReporters add reporters
func (handler *InitHandler) addReporters(c Community, r Repository, mapReporters, mapReportersInDB map[string]string) error {
	// reporters added
	listOfAddReporters := make([]string, 0)
	for rt := range mapReporters {
		_, okinReporters := mapReportersInDB[rt]
		if !okinReporters {
			listOfAddReporters = append(listOfAddReporters, rt)
		}
	}
	glog.Infof("list of add reporters: %v", listOfAddReporters)
	if len(listOfAddReporters) > 0 {
		// build create project member param
		memberbody := gitee.ProjectMemberPutParam{}
		memberbody.AccessToken = handler.Config.GiteeToken
		// memberbody.Permission = PermissionPull

		glog.Infof("begin to create reporters for: %s", *r.Name)
		for j := 0; j < len(listOfAddReporters); j++ {
			_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfAddReporters[j], memberbody)
			if err != nil {
				glog.Errorf("fail to create reporters: %v", err)
				continue
			}
			// create privilege
			ps := database.Privileges{
				Owner: *c.Name,
				Repo:  *r.Name,
				User:  listOfAddReporters[j],
				Type:  PrivilegeViewer,
			}
			err = database.DBConnection.Create(&ps).Error
			if err != nil {
				glog.Errorf("fail to create reporters in database: %v", err)
			}
		}
		glog.Infof("end to create reporters for: %s", *r.Name)
	}
	return nil
}

// removeReporters remove reporters
func (handler *InitHandler) removeReporters(c Community, r Repository, mapReporters, mapReportersInDB map[string]string) error {
	// reporters removed
	listOfRemoveReporters := make([]string, 0)
	for rt := range mapReportersInDB {
		_, okinReporters := mapReporters[rt]
		if !okinReporters {
			listOfRemoveReporters = append(listOfRemoveReporters, rt)
		}
	}
	glog.Infof("list of removed reporters: %v", listOfRemoveReporters)
	if len(listOfRemoveReporters) > 0 {
		// build remove project member param
		memberbody := &gitee.DeleteV5ReposOwnerRepoCollaboratorsUsernameOpts{}
		memberbody.AccessToken = optional.NewString(handler.Config.GiteeToken)

		glog.Infof("begin to remove reporters for: %s", *r.Name)
		for j := 0; j < len(listOfRemoveReporters); j++ {
			_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, *c.Name, *r.Name, listOfRemoveReporters[j], memberbody)
			if err != nil {
				glog.Errorf("fail to remove reporters: %v", err)
				continue
			}
			// delete privilege
			id, _ := strconv.Atoi(mapReportersInDB[listOfRemoveReporters[j]])
			ps := database.Privileges{}
			ps.ID = uint(id)
			err = database.DBConnection.Delete(&ps).Error
			if err != nil {
				glog.Errorf("fail to remove reporters in database: %v", err)
			}
		}
		glog.Infof("end to remove reporters for: %s", *r.Name)
	}
	return nil
}

// handleBranches handle branches
// currently for protecting branches only
func (handler *InitHandler) handleBranches(c Community, r Repository) error {
	// if the branches are defined in the repositories, it means that
	// all the branches defined in the community will not inherited by repositories
	mapBranches := make(map[string]string)

	if len(r.ProtectedBranches) > 0 {
		// using repository branches
		glog.Infof("using repository branches: %s", *r.Name)
		for _, b := range r.ProtectedBranches {
			mapBranches[b] = b
		}
	} else {
		// using community branches
		glog.Infof("using community branches: %s", *r.Name)
		for _, b := range c.ProtectedBranches {
			mapBranches[b] = b
		}
	}

	// get branches from DB
	var bs []database.Branches
	err := database.DBConnection.Model(&database.Branches{}).
		Where("owner = ? and repo = ?", c.Name, r.Name).Find(&bs).Error
	if err != nil {
		glog.Errorf("unable to get branches: %v", err)
		return err
	}
	mapBranchesInDB := make(map[string]string)
	for _, b := range bs {
		mapBranchesInDB[b.Name] = strconv.Itoa(int(b.ID))
	}

	// un-protected branches
	err = handler.removeBranchProtections(c, r, mapBranches, mapBranchesInDB)
	if err != nil {
		glog.Errorf("unable to un-protected branches: %v", err)
	}

	// protected branches
	err = handler.addBranchProtections(c, r, mapBranches, mapBranchesInDB)
	if err != nil {
		glog.Errorf("unable to protected branches: %v", err)
	}

	return nil
}

// unprotectedBranches unprotect branches
func (handler *InitHandler) removeBranchProtections(c Community, r Repository, mapBranches, mapBranchesInDB map[string]string) error {
	// remove branch protections
	listOfUnprotectedBranches := make([]string, 0)

	for k := range mapBranchesInDB {
		if _, exists := mapBranches[k]; !exists {
			listOfUnprotectedBranches = append(listOfUnprotectedBranches, k)
		}
	}
	glog.Infof("list of un-protected branches: %v", listOfUnprotectedBranches)

	if len(listOfUnprotectedBranches) > 0 {
		opts := &gitee.DeleteV5ReposOwnerRepoBranchesBranchProtectionOpts{}
		opts.AccessToken = optional.NewString(handler.Config.GiteeToken)

		glog.Infof("begin to remove branch protections for %s", *r.Name)
		for _, v := range listOfUnprotectedBranches {
			// remove branch protection from gitee
			_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoBranchesBranchProtection(
				handler.Context, *c.Name, *r.Name, v, opts)
			if err != nil {
				glog.Errorf("failed to remove branch protection: %v", err)
				continue
			}
			// remove branch protection from DB
			id, _ := strconv.Atoi(mapBranchesInDB[v])
			bs := database.Branches{}
			bs.ID = uint(id)
			err = database.DBConnection.Delete(&bs).Error
			if err != nil {
				glog.Errorf("failed to remove branch protection in database: %v", err)
			}
		}
		glog.Infof("end to remove branch protections for %s", *r.Name)
	}

	return nil
}

// addBranchProtections protects branches
func (handler *InitHandler) addBranchProtections(c Community, r Repository, mapBranches, mapBranchesInDB map[string]string) error {
	// add branch protections
	listOfProtectedBranches := make([]string, 0)

	for k := range mapBranches {
		if _, exits := mapBranchesInDB[k]; !exits {
			listOfProtectedBranches = append(listOfProtectedBranches, k)
		}
	}
	glog.Infof("list of protected branches: %v", listOfProtectedBranches)

	if len(listOfProtectedBranches) > 0 {
		getOpts := &gitee.GetV5ReposOwnerRepoBranchesBranchOpts{}
		getOpts.AccessToken = optional.NewString(handler.Config.GiteeToken)

		protectBody := gitee.BranchProtectionPutParam{}
		protectBody.AccessToken = handler.Config.GiteeToken

		glog.Infof("begin to add branch protections for %s", *r.Name)
		for _, v := range listOfProtectedBranches {
			// check if protected branch exists
			branchObj, response, _ := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoBranchesBranch(
				handler.Context, *c.Name, *r.Name, v, getOpts)
			if response.StatusCode == 404 {
				glog.Errorf("branch %s not exists, no need for protection", v)
				continue
			}

			// If branch has alreay been protected, no need for protection
			if branchObj.Protected == "true" {
				glog.Errorf("branch %s has been protected already, no need for protection", v)
				continue
			}

			// add branch protection to gitee
			_, response, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoBranchesBranchProtection(
				handler.Context, *c.Name, *r.Name, v, protectBody)
			if err != nil {
				glog.Errorf("failed to add branch protection: %v", err)
				continue
			}
			// add branch protection to database
			bs := database.Branches{
				Owner: *c.Name,
				Repo:  *r.Name,
				Name:  v,
				Type:  BranchProtected,
			}
			err = database.DBConnection.Create(&bs).Error
			if err != nil {
				glog.Errorf("failed to add branch protection in database: %v", err)
			}
		}
		glog.Infof("end to add branch protections for %s", *r.Name)
	}

	return nil
}

// handleRepositoryTypes handles that the repo is private or public
func (handler *InitHandler) handleRepositoryTypes(c Community, r Repository) error {
	// get repos from DB
	var rs database.Repositories
	err := database.DBConnection.Model(&database.Repositories{}).
		Where("owner = ? and repo = ?", c.Name, r.Name).First(&rs).Error
	if err != nil {
		glog.Errorf("unable to get repositories files: %v", err)
		return err
	}

	// the type is changed
	if rs.Type != *r.Type {
		// set value
		isSetPrivate := false
		if *r.Type == "private" {
			isSetPrivate = true
		}

		// invoke query repository
		glog.Infof("begin to query repository: %s", *r.Name)
		localVarOptionals := &gitee.GetV5ReposOwnerRepoOpts{}
		localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
		pj, response, _ := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepo(
			handler.Context, *c.Name, *r.Name, localVarOptionals)
		if response.StatusCode == 404 {
			glog.Infof("repository is not exist: %s", *r.Name)
			return nil
		}
		if pj.Private == isSetPrivate {
			glog.Infof("repository type is already: %s", *r.Type)
			return nil
		}

		// build patch repository param
		patchBody := gitee.RepoPatchParam{}
		patchBody.AccessToken = handler.Config.GiteeToken
		patchBody.Name = pj.Name
		patchBody.Description = pj.Description
		patchBody.Homepage = pj.Homepage
		patchBody.HasIssues = pj.HasIssues
		patchBody.HasWiki = pj.HasWiki
		patchBody.Description = pj.DefaultBranch
		patchBody.Private = isSetPrivate
		// invoke set type
		_, _, err = handler.GiteeClient.RepositoriesApi.PatchV5ReposOwnerRepo(handler.Context, *c.Name, *r.Name, patchBody)
		if err != nil {
			glog.Errorf("unable to set repository type: %v", err)
			return err
		}
	}

	return nil
}
