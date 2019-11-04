package cibot

import (
	"context"
	"encoding/base64"
	"time"

	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

type InitHandler struct {
	Config       config.Config
	Context      context.Context
	GiteeClient  *gitee.APIClient
	ProjectsFile string
}

type Projects struct {
	Community    Community    `yaml:"community"`
	Repostiories []Repostiory `yaml:"repostiories"`
}

type Community struct {
	Name      *string  `yaml:"name"`
	Manager   []string `yaml:"manager"`
	Developer []string `yaml:"developer"`
	Viewer    []string `yaml:"viewer"`
	Reporter  []string `yaml:"reporter"`
}

type Repostiory struct {
	Name        *string  `yaml:"name"`
	Description *string  `yaml:"description"`
	Type        *string  `yaml:"type"`
	Manager     []string `yaml:"manager"`
	Developer   []string `yaml:"developer"`
	Viewer      []string `yaml:"viewer"`
	Reporter    []string `yaml:"reporter"`
}

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
	// get params
	watchOwner := handler.Config.WatchProjectFileOwner
	watchRepo := handler.Config.WatchprojectFileRepo
	watchPath := handler.Config.WatchprojectFilePath
	watchRef := handler.Config.WatchProjectFileRef

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
	return nil
}

// watch database
func (handler *InitHandler) watch() {
	// get params
	watchOwner := handler.Config.WatchProjectFileOwner
	watchRepo := handler.Config.WatchprojectFileRepo
	watchPath := handler.Config.WatchprojectFilePath
	watchRef := handler.Config.WatchProjectFileRef
	watchDuration := handler.Config.WatchProjectFileDuration

	for {
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
									for i := 0; i < len(ps.Repostiories); i++ {
										// get repositories length
										lenRepositories, err := handler.getRepositoriesLength(*ps.Community.Name, *ps.Repostiories[i].Name, pf.ID)
										if err != nil {
											glog.Errorf("failed to get repositories length: %v", err)
											continue
										}
										if lenRepositories > 0 {
											glog.Infof("repository: %s is exist. no action.", *ps.Repostiories[i].Name)
										} else {
											// add repository
											err = handler.addRepositories(*ps.Community.Name, *ps.Repostiories[i].Name,
												*ps.Repostiories[i].Description, *ps.Repostiories[i].Type, pf.ID)
											if err != nil {
												glog.Errorf("failed to add repositories: %v", err)
											}
										}
									}
								}
							}
						}
					}
				} else {
					glog.Infof("no waiting sha: %v", pf.WaitingSha)
				}
			}
		}

		// watch duration
		glog.Info("end to serve")
		time.Sleep(time.Duration(watchDuration) * time.Second)
	}
}

// GetRepositoriesLength get repositories length
func (handler *InitHandler) getRepositoriesLength(owner string, repo string, id uint) (int, error) {
	// Check repositories file
	var lenRepositories int
	err := database.DBConnection.Model(&database.Repositories{}).
		Where("owner = ? and repo = ? and project_file_id = ?", owner, repo, id).
		Count(&lenRepositories).Error
	if err != nil {
		glog.Errorf("unable to get repositories files: %v", err)
	}
	return lenRepositories, err
}

// addRepositories add repository
func (handler *InitHandler) addRepositories(owner, repo, description, t string, id uint) error {
	// add repository in gitee
	err := handler.addRepositoriesinGitee(owner, repo, description, t)
	if err != nil {
		glog.Errorf("failed to add repositories: %v", err)
		return err
	}

	// add repository in database
	err = handler.addRepositoriesinDB(owner, repo, description, t, id)
	if err != nil {
		glog.Errorf("failed to add repositories: %v", err)
		return err
	}
	return nil
}

// addRepositoriesinDB add repository in database
func (handler *InitHandler) addRepositoriesinDB(owner, repo, description, t string, id uint) error {
	// add repository
	addrepo := database.Repositories{
		Owner:         owner,
		Repo:          repo,
		Description:   description,
		Type:          t,
		ProjectFileID: id,
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
	if t == "private" {
		repobody.Private = true
	} else {
		repobody.Private = false
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
