package cibot

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

type OwnerHandler struct {
	Config      config.Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

// Serve
func (handler *OwnerHandler) Serve() {
	// watch database
	handler.watch()
}

// watch database
func (handler *OwnerHandler) watch() {
	for {
		watchDuration := handler.Config.WatchOwnerFileDuration
		// get repositories from DB
		var rs []database.Repositories
		err := database.DBConnection.Model(&database.Repositories{}).Find(&rs).Error
		if err != nil {
			glog.Errorf("unable to get repos: %v", err)
		} else {
			if len(rs) > 0 {
				// get sigs from DB
				var srs []database.SigRecords
				err := database.DBConnection.Model(&database.SigRecords{}).Find(&srs).Error
				if err != nil {
					glog.Errorf("unable to get sigs: %v", err)
				} else {
					if len(srs) > 0 {
						// get owner from files
						getOwnersResult := true
						mapSigOwners := make(map[string]map[string]string)
						for _, sr := range srs {
							for _, wf := range handler.Config.WatchOwnerFiles {
								mapOwners := make(map[string]string)
								owners, err := handler.getOwners(wf, sr.Name)
								if err != nil {
									// set result into false
									getOwnersResult = false
									glog.Errorf("unable to getOwners: %v", err)
								}
								if len(owners) > 0 {
									for _, o := range owners {
										mapOwners[o] = o
									}
								}
								mapSigOwners[sr.Name] = mapOwners
								break
							}
						}

						// based on repository
						if getOwnersResult {
							for _, repo := range rs {
								err = handler.handleOwners(repo, mapSigOwners)
								if err != nil {
									glog.Errorf("unable to handle owners: %v", err)
								}
							}
						}
					} else {
						glog.Info("current sig length is zero")
					}
				}
			} else {
				glog.Info("current repository length is zero")
			}
		}

		// watch duration
		glog.Info("end to serve in owner")
		time.Sleep(time.Duration(watchDuration) * time.Second)
	}
}

// getOwners get owners
func (handler *OwnerHandler) getOwners(wf config.WatchOwnerFile, name string) ([]string, error) {
	// get params
	watchOwner := wf.WatchOwnerFileOwner
	watchRepo := wf.WatchOwnerFileRepo
	watchPath := strings.Replace(wf.WatchOwnerFilePath, "*", name, -1)
	watchRef := wf.WatchOwnerFileRef

	// invoke api to get file contents
	localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
	localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
	localVarOptionals.Ref = optional.NewString(watchRef)

	// get contents
	contents, _, err := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
		handler.Context, watchOwner, watchRepo, watchPath, localVarOptionals)
	if err != nil {
		glog.Errorf("unable to get repository content: %v", err)
		return nil, err
	}

	// base64 decode
	decodeBytes, err := base64.StdEncoding.DecodeString(contents.Content)
	if err != nil {
		glog.Errorf("decode content with error: %v", err)
		return nil, err
	}
	// unmarshal owners file
	var owners OwnersFile
	err = yaml.Unmarshal(decodeBytes, &owners)
	if err != nil {
		glog.Errorf("fail to unmarshal owners: %v", err)
		return nil, err
	}

	// return owners
	if len(owners.Maintainers) > 0 {
		return owners.Maintainers, nil
	}

	return nil, nil
}

// handleOwners handle owners
func (handler *OwnerHandler) handleOwners(repo database.Repositories, mapSigOwners map[string]map[string]string) error {
	// get owner for sig
	var srepos []database.SigRepositories
	err := database.DBConnection.Model(&database.SigRepositories{}).
		Where("repo_name = ?", repo.Owner+"/"+repo.Repo).Find(&srepos).Error
	if err != nil {
		glog.Errorf("unable to get sig repos: %v", err)
	}
	expectedMembers := make(map[string]string)
	for _, srepo := range srepos {
		for k, v := range mapSigOwners[srepo.Name] {
			expectedMembers[k] = v
		}
	}

	// get current owners
	var ps []database.Privileges
	err = database.DBConnection.Model(&database.Privileges{}).
		Where("owner = ? and repo = ? and type = ?", repo.Owner, repo.Repo, PrivilegeDeveloper).Find(&ps).Error
	if err != nil {
		glog.Errorf("unable to get members: %v", err)
	}
	actualMembers := make(map[string]string)
	if len(ps) > 0 {
		for _, p := range ps {
			actualMembers[p.User] = strconv.Itoa(int(p.ID))
		}
	}

	// remove
	err = handler.removeOwners(repo, expectedMembers, actualMembers)
	if err != nil {
		glog.Errorf("unable to remove sig: %v", err)
	}

	// add
	err = handler.addOwners(repo, expectedMembers, actualMembers)
	if err != nil {
		glog.Errorf("unable to add sig: %v", err)
	}

	return nil
}

// removeOwners
func (handler *OwnerHandler) removeOwners(repo database.Repositories, expectedMembers, actualMembers map[string]string) error {
	listOfRemove := make([]string, 0)

	for k := range actualMembers {
		if _, exists := expectedMembers[k]; !exists {
			listOfRemove = append(listOfRemove, k)
		}
	}

	if len(listOfRemove) > 0 {
		glog.Infof("list of remove privileges: %v", listOfRemove)
		memberbody := &gitee.DeleteV5ReposOwnerRepoCollaboratorsUsernameOpts{}
		memberbody.AccessToken = optional.NewString(handler.Config.GiteeToken)

		glog.Infof("begin to remove privileges for %s/%s", repo.Owner, repo.Repo)
		for _, v := range listOfRemove {
			_, err := handler.GiteeClient.RepositoriesApi.DeleteV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, repo.Owner, repo.Repo, v, memberbody)
			if err != nil {
				glog.Errorf("fail to remove privileges: %v", err)
				continue
			}

			// remove from DB
			id, _ := strconv.Atoi(actualMembers[v])
			sr := database.Privileges{}
			sr.ID = uint(id)
			err = database.DBConnection.Delete(&sr).Error
			if err != nil {
				glog.Errorf("failed to remove privilege in database: %v", err)
			}
		}
		glog.Infof("end to remove privileges for %s/%s", repo.Owner, repo.Repo)
	}

	return nil
}

// addOwners
func (handler *OwnerHandler) addOwners(repo database.Repositories, expectedMembers, actualMembers map[string]string) error {
	listOfAdd := make([]string, 0)

	for k := range expectedMembers {
		if _, exits := actualMembers[k]; !exits {
			listOfAdd = append(listOfAdd, k)
		}
	}

	if len(listOfAdd) > 0 {
		glog.Infof("list of add privileges: %v", listOfAdd)
		memberbody := gitee.ProjectMemberPutParam{}
		memberbody.AccessToken = handler.Config.GiteeToken
		memberbody.Permission = PermissionPush

		glog.Infof("begin to add privileges for %s/%s", repo.Owner, repo.Repo)
		for _, v := range listOfAdd {
			_, _, err := handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(
				handler.Context, repo.Owner, repo.Repo, v, memberbody)
			if err != nil {
				glog.Errorf("fail to create developers: %v", err)
				continue
			}
			// create privilege
			ps := database.Privileges{
				Owner: repo.Owner,
				Repo:  repo.Repo,
				User:  v,
				Type:  PrivilegeDeveloper,
			}
			err = database.DBConnection.Create(&ps).Error
			if err != nil {
				glog.Errorf("failed to add privileges in database: %v", err)
			}
		}
		glog.Infof("end to add privileges for %s/%s", repo.Owner, repo.Repo)
	}

	return nil
}
