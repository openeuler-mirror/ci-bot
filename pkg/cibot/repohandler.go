package cibot

import (
	"context"
	"encoding/base64"
	"fmt"
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
	Version      string       `yaml:"version"`
	Community    string       `yaml:"community"`
	Repositories []Repository `yaml:"repositories"`
}

type Branch struct {
	Name       *string `yaml:"name"`
	Type       *string `yaml:"type"`
	CreateFrom *string `yaml:"create_from"`
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
			glog.Infof("project file exists: %s", contents.Sha)
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
			glog.Infof("project file does not exist: %s", contents.Sha)
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
											lenRepositories, errex := handler.getRepositoriesLength(ps.Community, *ps.Repositories[i].Name)
											if errex != nil {
												glog.Errorf("failed to get repositories length: %v", errex)
												result = false
												continue
											}
											if lenRepositories > 0 {
												glog.Infof("repository: %s exists. no action.", *ps.Repositories[i].Name)
											} else {
												// add repository
												err = handler.addRepositories(ps.Community, ps.Repositories[i])
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
											// handle repository settings, currently type and can_comment are supported
											err = handler.handleRepositorySetting(ps.Community, ps.Repositories[i])
											if err != nil {
												glog.Errorf("failed to handle repository setting: %v", err)
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
func (handler *RepoHandler) addRepositories(owner string, repo Repository) error {
	// add repository in gitee
	err := handler.addRepositoriesinGitee(owner, repo)
	if err != nil {
		glog.Errorf("failed to add repositories: %v", err)
		return err
	}

	// add repository in database
	err = handler.addRepositoriesinDB(owner, repo)
	if err != nil {
		glog.Errorf("failed to add repositories: %v", err)
		return err
	}
	return nil
}

// addRepositoriesinDB add repository in database
func (handler *RepoHandler) addRepositoriesinDB(owner string, repo Repository) error {
	// this is a rename instead of create operation,
	// so only update existing repository in DB
	if repo.RenameFrom != nil && len(*repo.RenameFrom) > 0 {
		// Update the repositories and branches table in a single transactions
		tx := database.DBConnection.Begin()
		err := tx.Model(&database.Repositories{}).
			Where("owner = ? and repo = ?", owner, *repo.RenameFrom).
			Update("repo", *repo.Name).Error
		if err != nil {
			glog.Errorf("unable to rename repository: %v", err)
			tx.Rollback()
			return err
		}
		err = tx.Model(&database.Branches{}).
			Where("owner = ? and repo = ?", owner, *repo.RenameFrom).
			Update("repo", *repo.Name).Error
		if err != nil {
			glog.Errorf("unable to rename repository for relevant branches: %v", err)
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	// add repository
	addrepo := database.Repositories{
		Owner:       owner,
		Repo:        *repo.Name,
		Description: *repo.Description,
		Type:        *repo.Type,
		Commentable: repo.IsCommentable(),
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
func (handler *RepoHandler) addRepositoriesinGitee(owner string, repo Repository) error {
	// build create repository param
	repobody := gitee.RepositoryPostParam{}
	repobody.AccessToken = handler.Config.GiteeToken
	repobody.Name = *repo.Name
	repobody.Description = *repo.Description
	repobody.HasIssues = true
	repobody.HasWiki = true
	// set `auto_init` as true to initialize `master` branch with README after repo creation
	repobody.AutoInit = true
	repobody.CanComment = repo.IsCommentable()
	repobody.Private = *repo.Type == "private"

	// invoke query repository
	glog.Infof("begin to query repository: %s", *repo.Name)
	localVarOptionals := &gitee.GetV5ReposOwnerRepoOpts{}
	localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
	_, response, _ := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepo(handler.Context, owner, *repo.Name, localVarOptionals)
	if response.StatusCode == 404 {
		glog.Infof("repository does not exist: %s", *repo.Name)
	} else {
		glog.Infof("repository has already existed: %s", *repo.Name)
		return nil
	}

	// if rename_from does exist, now go to invoke rename repository
	if repo.RenameFrom != nil && len(*repo.RenameFrom) > 0 {
		// invoke query repoisitory with the name defined by rename_from
		glog.Infof("begin to query repository: %s defined by rename_from ", *repo.RenameFrom)
		localVarRenameOptionals := &gitee.GetV5ReposOwnerRepoOpts{}
		localVarRenameOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
		_, response, _ = handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepo(handler.Context, owner, *repo.RenameFrom, localVarRenameOptionals)
		if response.StatusCode == 404 {
			errMsg := fmt.Sprintf("repository defined by rename_from does not exist: %s", *repo.RenameFrom)
			glog.Errorf("failed to rename repository: %s", errMsg)
			return fmt.Errorf(errMsg)
		}
		// everything seems fine, then build patch repository param
		repoPatchParam := gitee.RepoPatchParam{}
		repoPatchParam.AccessToken = handler.Config.GiteeToken
		repoPatchParam.Name = *repo.Name
		repoPatchParam.Path = *repo.Name
		// invoke patching repository API to change *repo.RenameFrom to *repo.Name
		_, _, err := handler.GiteeClient.RepositoriesApi.PatchV5ReposOwnerRepo(handler.Context, owner, *repo.RenameFrom, repoPatchParam)
		if err != nil {
			glog.Errorf("unable to rename the repository from %s to %s", *repo.RenameFrom, *repo.Name)
			return err
		}
		return nil
	}

	// invoke create repository
	glog.Infof("begin to create repository: %s", *repo.Name)
	_, _, err := handler.GiteeClient.RepositoriesApi.PostV5OrgsOrgRepos(handler.Context, owner, repobody)
	if err != nil {
		glog.Errorf("fail to create repository: %v", err)
		return err
	}
	glog.Infof("end to create repository: %s", *repo.Name)

	// create branch
	repobranchbody := gitee.CreateBranchParam{}
	repobranchbody.AccessToken = handler.Config.GiteeToken
	for _, br := range repo.Branches {
		if *br.Name == "master" {
			continue
		}
		repobranchbody.BranchName = *br.Name
		if *br.CreateFrom == "" {
			*br.CreateFrom = "master"
		}
		repobranchbody.Refs = *br.CreateFrom
		_, _, err := handler.GiteeClient.RepositoriesApi.PostV5ReposOwnerRepoBranches(handler.Context, owner, *repo.Name, repobranchbody)
		if err != nil {
			glog.Errorf("fail to add branch (%s) for repository (%s): %v", *br.Name, *repo.Name, err)
			return err
		}
		glog.Infof("Add branch (%s) for repository (%s)", *br.Name, *repo.Name)
	}
	return nil
}

// handleBranches handle branches
// currently for protecting branches only
func (handler *RepoHandler) handleBranches(community string, r Repository) error {
	// if the branches are defined in the repositories, it means that
	// all the branches defined in the community will not inherited by repositories

	if len(r.Branches) > 0 {
		// using repository branches
		glog.Infof("Setting repository branches: %s", *r.Name)
		for _, b := range r.Branches {
			//check yaml branches is in db
			var bs []database.Branches
			err := database.DBConnection.Model(&database.Branches{}).
				Where("owner = ? and repo = ? and name = ?", community, r.Name, b.Name).Find(&bs).Error
			if err != nil || len(bs) == 0 {
				glog.Infof("Do not have branches (%s), need to create. %v, ", *b.Name, err)
				err = handler.addBranchGiteeAndDb(community, r, b)
				if err != nil {
					glog.Errorf("Add branch(%s) of repor(%s) failed: %v", *b.Name, *r.Name, err)
				}
				continue
			}

			glog.Infof("Branch exist, branch(%s) of repo(%s).", *b.Name, *r.Name)
			err = handler.changeBranchGiteeAndDb(community, r, b, bs[0])
			if err != nil {
				glog.Errorf("Change branch(%s) of repor(%s) failed: %v", *b.Name, *r.Name, err)
			}
			continue
		}
	}
	return nil
}

// unprotectedBranches unprotect branches
func (handler *RepoHandler) changeBranchGiteeAndDb(community string, r Repository, br Branch, bs database.Branches) error {
	// if branch features are same as ones in db, do nothing.
	if *r.Name == bs.Repo && *br.Name == bs.Name && *br.Type == bs.Type {
		return nil
	}

	// change branch protected freature in Gitee
	var brType string
	brType = BranchNormal
	if *br.Type == BranchProtected {
		protectBody := gitee.BranchProtectionPutParam{}
		protectBody.AccessToken = handler.Config.GiteeToken
		_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoBranchesBranchProtection(
			handler.Context, community, *r.Name, *br.Name, protectBody)
		if err != nil {
			glog.Errorf("failed to add branch protection: %v", err)
		}
		brType = BranchProtected
	} else {
		opts := &gitee.DeleteV5ReposOwnerRepoBranchesBranchProtectionOpts{}
		opts.AccessToken = optional.NewString(handler.Config.GiteeToken)
		_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoBranchesBranchProtection(
			handler.Context, community, *r.Name, *br.Name, opts)
		if err != nil {
			glog.Errorf("failed to remove branch protection: %v", err)
		}
	}

	// change branch protected freature in DB
	updatebranch := &database.Branches{}
	updatebranch.ID = bs.ID
	err := database.DBConnection.Model(updatebranch).Update("Type", brType).Error
	if err != nil {
		glog.Errorf("unable to update type: %v", err)
	}

	return nil
}

// addBranchProtections protects branches
func (handler *RepoHandler) addBranchGiteeAndDb(community string, r Repository, br Branch) error {
	// create branch in gitee
	repobranchbody := gitee.CreateBranchParam{}
	repobranchbody.AccessToken = handler.Config.GiteeToken
	if br.Name != nil {
		repobranchbody.BranchName = *br.Name
	} else {
		return nil
	}
	if br.CreateFrom != nil {
		repobranchbody.Refs = *br.CreateFrom
	} else {
		repobranchbody.Refs = ""
	}

	_, _, err := handler.GiteeClient.RepositoriesApi.PostV5ReposOwnerRepoBranches(handler.Context, community, *r.Name, repobranchbody)
	if err != nil {
		glog.Errorf("fail to add branch (%s) for repository (%s): %v", *br.Name, *r.Name, err)
	} else {
                glog.Infof("Add branch (%s) for repository (%s) in gitee.", *br.Name, *r.Name)
	}

	// add branch protection to gitee
	var brType string
	brType = BranchNormal
	if *br.Type == BranchProtected {
		protectBody := gitee.BranchProtectionPutParam{}
		protectBody.AccessToken = handler.Config.GiteeToken
		_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoBranchesBranchProtection(
			handler.Context, community, *r.Name, *br.Name, protectBody)
		if err != nil {
			glog.Errorf("failed to add branch protection: %v", err)
		}
		brType = BranchProtected
	}
	// add branch to database
	bs := database.Branches{
		Owner:          community,
		Repo:           *r.Name,
		Name:           *br.Name,
		Type:           brType,
		AdditionalInfo: repobranchbody.Refs,
	}
	err = database.DBConnection.Create(&bs).Error
	if err != nil {
		glog.Errorf("failed to add branch protection in database: %v", err)
	}

	return nil
}

// handleRepositorySetting handles that the repo settings, including:
//   1. type is private or public
//   2. commentable is true or false
//   3. set none reviewer but not ci-bot(default)
func (handler *RepoHandler) handleRepositorySetting(community string, r Repository) error {
	// get repos from DB
	var rs database.Repositories
	err := database.DBConnection.Model(&database.Repositories{}).
		Where("owner = ? and repo = ?", community, r.Name).First(&rs).Error
	if err != nil {
		glog.Errorf("unable to get repositories: %v", err)
		return err
	}

	typePrivateExpected := (*r.Type == "private")
	// mark type as changed if the type from DB is NOT identical to yaml
	typeChanged := (rs.Type != *r.Type)

	commentableExpected := r.IsCommentable()
	// mark commentable as changed if the commentable from DB is NOT identical to yaml
	commentableChanged := (rs.Commentable != commentableExpected)

	if typeChanged || commentableChanged {
		// invoke query repository
		glog.Infof("begin to query repository: %s", *r.Name)
		localVarOptionals := &gitee.GetV5ReposOwnerRepoOpts{}
		localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
		pj, response, _ := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepo(
			handler.Context, community, *r.Name, localVarOptionals)
		if response.StatusCode == 404 {
			glog.Infof("repository dose not exist: %s", *r.Name)
			return nil
		}

		// do we need to invoke the gitee patch api to change the setting of the repository
		needInvokeRepoPatchAPI := false
		// handle the change of type
		if typeChanged {
			// if the repo type from gitee is identical to yaml
			if pj.Private == typePrivateExpected {
				glog.Infof("repository type is already: %s", *r.Type)
			} else {
				needInvokeRepoPatchAPI = true
				glog.Infof("going to change repo type via gitee API")
			}
		}
		// handle the change of commentable
		if commentableChanged {
			// if the repo commentable from gitee is identical to yaml
			if pj.CanComment == commentableExpected {
				glog.Infof("repository commentable is already: %t", pj.CanComment)
			} else {
				needInvokeRepoPatchAPI = true
				glog.Infof("going to change repo commentable via gitee API")
			}
		}

		// now to invoke gitee api to change the repo settings
		if needInvokeRepoPatchAPI {
			// build patch repository param
			patchBody := gitee.RepoPatchParam{}
			patchBody.AccessToken = handler.Config.GiteeToken
			patchBody.Name = pj.Name
			patchBody.Description = pj.Description
			patchBody.Homepage = pj.Homepage
			patchBody.HasIssues = strconv.FormatBool(pj.HasIssues)
			patchBody.HasWiki = strconv.FormatBool(pj.HasWiki)
			patchBody.Private = strconv.FormatBool(typePrivateExpected)
			patchBody.CanComment = strconv.FormatBool(commentableExpected)

			// invoke set type
			_, _, err = handler.GiteeClient.RepositoriesApi.PatchV5ReposOwnerRepo(handler.Context, community, *r.Name, patchBody)
			if err != nil {
				glog.Errorf("unable to set repository settings: %v", err)
				return err
			}
		}

		// define update repository
		updaterepo := &database.Repositories{}
		updaterepo.ID = rs.ID
		err = database.DBConnection.Model(updaterepo).
			Update("Type", *r.Type).
			Update("Commentable", commentableExpected).
			Error
		if err != nil {
			glog.Errorf("unable to update repository settings: %v", err)
		}
	}

	// set none reviewer but not ci-bot(default)

	//	reviewerBody := gitee.SetRepoReviewer{}
	//	reviewerBody.AccessToken = handler.Config.GiteeToken
	//	reviewerBody.Assignees = " "
	//	reviewerBody.Testers = " "
	//	reviewerBody.AssigneesNumber = 0
	//	reviewerBody.TestersNumber = 0
	//	response, errex := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoReviewer(handler.Context, community, *r.Name, reviewerBody)
	//	if errex != nil {
	//		glog.Errorf("Set repository reviewer info failed: %v, %s", errex, response.Status)
	//		glog.Errorf("requestURL:%s, %s", response.Request.RequestURI, response.Request.Host)
	//		return errex
	//	}

	return nil
}
