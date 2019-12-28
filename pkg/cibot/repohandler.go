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

type RepoHandler struct {
	Config      config.Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

type Repos struct {
	Community    string       `yaml:"community"`
	Repositories []Repository `yaml:"repositories"`
}

// Serve
func (handler *RepoHandler) Serve() {
	// init sha
	err := handler.initSha()
	if err != nil {
		glog.Errorf("unable to initSha: %v", err)
		return
	}
	// watch database
	handler.watch()
}

// initSha init sha
func (handler *RepoHandler) initSha() error {
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
			// init targetsha
			updatepf.TargetSha = ""
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
func (handler *RepoHandler) watch() {
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
									var ps Repos
									err = yaml.Unmarshal(decodeBytes, &ps)
									if err != nil {
										glog.Errorf("failed to unmarshal repos: %v", err)
									} else {
										glog.Infof("get blob result: %v", ps)
										result := true
										for i := 0; i < len(ps.Repositories); i++ {
											// get repositories length
											lenRepositories, err := handler.getRepositoriesLength(ps.Community, *ps.Repositories[i].Name)
											if err != nil {
												glog.Errorf("failed to get repositories length: %v", err)
												result = false
												continue
											}
											if lenRepositories > 0 {
												glog.Infof("repository: %s is exist. no action.", *ps.Repositories[i].Name)
											} else {
												// add repository
												err = handler.addRepositories(ps.Community, *ps.Repositories[i].Name,
													*ps.Repositories[i].Description, *ps.Repositories[i].Type)
												if err != nil {
													glog.Errorf("failed to add repositories: %v", err)
													result = false
													continue
												}
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
func (handler *RepoHandler) getRepositoriesLength(owner string, repo string) (int, error) {
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
func (handler *RepoHandler) addRepositories(owner, repo, description, t string) error {
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
func (handler *RepoHandler) addRepositoriesinDB(owner, repo, description, t string) error {
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
func (handler *RepoHandler) addRepositoriesinGitee(owner, repo, description, t string) error {
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

// handleBranches handle branches
// currently for protecting branches only
func (handler *RepoHandler) handleBranches(community string, r Repository) error {
	// if the branches are defined in the repositories, it means that
	// all the branches defined in the community will not inherited by repositories
	mapBranches := make(map[string]string)

	if len(r.ProtectedBranches) > 0 {
		// using repository branches
		glog.Infof("using repository branches: %s", *r.Name)
		for _, b := range r.ProtectedBranches {
			mapBranches[b] = b
		}
	}

	// get branches from DB
	var bs []database.Branches
	err := database.DBConnection.Model(&database.Branches{}).
		Where("owner = ? and repo = ?", community, r.Name).Find(&bs).Error
	if err != nil {
		glog.Errorf("unable to get branches: %v", err)
		return err
	}
	mapBranchesInDB := make(map[string]string)
	for _, b := range bs {
		mapBranchesInDB[b.Name] = strconv.Itoa(int(b.ID))
	}

	// un-protected branches
	err = handler.removeBranchProtections(community, r, mapBranches, mapBranchesInDB)
	if err != nil {
		glog.Errorf("unable to un-protected branches: %v", err)
	}

	// protected branches
	err = handler.addBranchProtections(community, r, mapBranches, mapBranchesInDB)
	if err != nil {
		glog.Errorf("unable to protected branches: %v", err)
	}

	return nil
}

// unprotectedBranches unprotect branches
func (handler *RepoHandler) removeBranchProtections(community string, r Repository, mapBranches, mapBranchesInDB map[string]string) error {
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
				handler.Context, community, *r.Name, v, opts)
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
func (handler *RepoHandler) addBranchProtections(community string, r Repository, mapBranches, mapBranchesInDB map[string]string) error {
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
				handler.Context, community, *r.Name, v, getOpts)
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
				handler.Context, community, *r.Name, v, protectBody)
			if err != nil {
				glog.Errorf("failed to add branch protection: %v", err)
				continue
			}
			// add branch protection to database
			bs := database.Branches{
				Owner: community,
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
func (handler *RepoHandler) handleRepositoryTypes(community string, r Repository) error {
	// get repos from DB
	var rs database.Repositories
	err := database.DBConnection.Model(&database.Repositories{}).
		Where("owner = ? and repo = ?", community, r.Name).First(&rs).Error
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
			handler.Context, community, *r.Name, localVarOptionals)
		if response.StatusCode == 404 {
			glog.Infof("repository is not exist: %s", *r.Name)
			return nil
		}
		if pj.Private == isSetPrivate {
			glog.Infof("repository type is already: %s", *r.Type)
		} else {
			// build patch repository param
			patchBody := gitee.RepoPatchParam{}
			patchBody.AccessToken = handler.Config.GiteeToken
			patchBody.Name = pj.Name
			patchBody.Description = pj.Description
			patchBody.Homepage = pj.Homepage
			if pj.HasIssues {
				patchBody.HasIssues = "true"
			} else {
				patchBody.HasIssues = "false"
			}
			if pj.HasWiki {
				patchBody.HasWiki = "true"
			} else {
				patchBody.HasWiki = "false"
			}
			if isSetPrivate {
				patchBody.Private = "true"
			} else {
				patchBody.Private = "false"
			}

			// invoke set type
			_, _, err = handler.GiteeClient.RepositoriesApi.PatchV5ReposOwnerRepo(handler.Context, community, *r.Name, patchBody)
			if err != nil {
				glog.Errorf("unable to set repository type: %v", err)
				return err
			}
		}

		// define update repository
		updaterepo := &database.Repositories{}
		updaterepo.ID = rs.ID
		err = database.DBConnection.Model(updaterepo).Update("Type", *r.Type).Error
		if err != nil {
			glog.Errorf("unable to update type: %v", err)
		}
	}

	return nil
}
