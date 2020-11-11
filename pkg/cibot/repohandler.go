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
	repobranchbody.Refs = "master"
	for _, br := range repo.ProtectedBranches {
		repobranchbody.BranchName = br
		_, _, err := handler.GiteeClient.RepositoriesApi.PostV5ReposOwnerRepoBranches(handler.Context, owner, *repo.Name, repobranchbody)
		if br == "master" {
			continue
		}
		if err != nil {
			glog.Errorf("fail to add branch (%s) for repository (%s): %v", br, *repo.Name, err)
			return err
		}
		glog.Infof("Add branch (%s) for repository (%s)", br, *repo.Name)
	}
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
				glog.Errorf("branch %s does not exist, no need for protection", v)
				continue
			}

			// If branch has alreay been protected, no need for protection
			if branchObj.Protected == true {
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
	reviewerBody := gitee.SetRepoReviewer{}
	reviewerBody.AccessToken = handler.Config.GiteeToken
	reviewerBody.Assignees = " "
	reviewerBody.Testers = " "
	reviewerBody.AssigneesNumber = 0
	reviewerBody.TestersNumber = 0
	response, errex := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoReviewer(handler.Context, community, *r.Name, reviewerBody)
	if errex != nil {
		glog.Errorf("Set repository reviewer info failed: %v, %s", errex, response.Status)
		return errex
	}

	return nil
}
