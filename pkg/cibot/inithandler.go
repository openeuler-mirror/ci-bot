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
	Community Community    `yaml:"community"`
	Projects  []Repostiory `yaml:"repostiories"`
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
	Manager     []string `yaml:"manager"`
	Developer   []string `yaml:"developer"`
	Viewer      []string `yaml:"viewer"`
	Reporter    []string `yaml:"reporter"`
	Type        *string  `yaml:"type"`
	Description *string  `yaml:"description"`
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
