package cibot

import (
	"strings"

	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
)

// HandlePushEvent handles push event
func (s *Server) HandlePushEvent(event *gitee.PushEvent) {
	if event == nil {
		return
	}

	// handle events
	if (event.Repository.Namespace == s.Config.WatchProjectFileOwner) && (event.Repository.Name == s.Config.WatchprojectFileRepo) {
		// owner and repo are matched
		if event.Ref != nil {
			ref := event.Ref
			configRef := s.Config.watchProjectFileRef
			if configRef == "" {
				configRef = "master"
			}
			// refs/heads/master
			if strings.Index(ref, configRef) >= 0 {
				// the branch is matched
				glog.Infof("push event is triggered. owner: %s repo: %s ref: %s", event.Repository.Namespace, event.Repository.Name, ref)
				// invoke api to get file contents
				localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
				localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)
				localVarOptionals.Ref = optional.NewString(configRef)
				// get contents
				contents, _, err := s.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
					s.Context, event.Repository.Namespace, event.Repository.Name, s.Config.WatchprojectFilePath, localVarOptionals)
				if err != nil {
					glog.Errorf("unable to get repository content by path: %v", err)
					return
				}
				glog.Infof("get triggered sha: %s", contents.Sha)

				// Check project file in database
				var lenProjectFiles int
				err = database.DBConnection.Model(&database.ProjectFiles{}).
					Where("owner = ? and repo = ? and path = ? and ref = ?",
						event.Repository.Namespace, event.Repository.Name, s.Config.WatchprojectFilePath, configRef).
					Count(&lenProjectFiles).Error
				if err != nil {
					glog.Errorf("unable to get project files in database: %v", err)
					return
				}
				if lenProjectFiles > 0 {
					glog.Infof("project file is exist. triggered sha: %s", contents.Sha)
					// Check sha in database
					updatepf := database.ProjectFiles{}
					err = database.DBConnection.
						Where("owner = ? and repo = ? and path = ? and ref = ?",
							event.Repository.Namespace, event.Repository.Name, s.Config.WatchprojectFilePath, configRef).
						First(&updatepf).Error
					if err != nil {
						glog.Errorf("unable to get project files in database: %v", err)
						return
					}
					glog.Infof("project file current sha: %v target sha: %v waiting sha: %v", updatepf.CurrentSha, updatepf.TargetSha, updatepf.WaitingSha)
					if (updatepf.CurrentSha != contents.Sha) && (updatepf.TargetSha != contents.Sha) && (updatepf.WaitingSha != contents.Sha) {
						// write sha in waitingsha
						updatepf.WaitingSha = contents.Sha
						err = database.DBConnection.Save(&updatepf).Error
						if err != nil {
							glog.Errorf("unable to save project files in database: %v", err)
							return
						}
						glog.Infof("update waiting sha successfully. triggered sha: %s", contents.Sha)
					}
				} else {
					glog.Infof("project file is non-exist. triggered sha: %s", contents.Sha)
					// add project file
					addpf := database.ProjectFiles{
						Owner:      event.Repository.Namespace,
						Repo:       event.Repository.Name,
						Path:       s.Config.WatchprojectFilePath,
						Ref:        configRef,
						WaitingSha: contents.Sha,
					}

					// create project file
					err = database.DBConnection.Create(&addpf).Error
					if err != nil {
						glog.Errorf("unable to create project files in database: %v", err)
						return
					}
					glog.Infof("add project file successfully. triggered sha: %s", contents.Sha)
				}
			}
		}
	}
}
