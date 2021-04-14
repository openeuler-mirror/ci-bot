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

type SigHandler struct {
	Config      config.Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

type SigsYaml struct {
	Sigs []Sig `yaml:"sigs"`
}

type Sig struct {
	Name         string   `yaml:"name"`
	Repositories []string `yaml:"repositories"`
}

// Serve
func (handler *SigHandler) Serve() {
	// init sha
	err := handler.initSha()
	if err != nil {
		glog.Errorf("unable to initSha in sig: %v", err)
		return
	}
	// watch database
	handler.watch()
}

// initSha init sha
func (handler *SigHandler) initSha() error {
	if len(handler.Config.WatchSigFiles) == 0 {
		return nil
	}

	for _, wf := range handler.Config.WatchSigFiles {
		// get params
		watchOwner := wf.WatchSigFileOwner
		watchRepo := wf.WatchSigFileRepo
		watchPath := wf.WatchSigFilePath
		watchRef := wf.WatchSigFileRef

		// invoke api to get file contents
		localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
		localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
		localVarOptionals.Ref = optional.NewString(watchRef)

		// get contents
		contents, _, err := handler.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
			handler.Context, watchOwner, watchRepo, watchPath, localVarOptionals)
		if err != nil {
			glog.Errorf("unable to get repository content in sig: %v", err)
			return err
		}
		// Check sig file
		var lenSigFiles int
		err = database.DBConnection.Model(&database.SigFiles{}).
			Where("owner = ? and repo = ? and path = ? and ref = ?", watchOwner, watchRepo, watchPath, watchRef).
			Count(&lenSigFiles).Error
		if err != nil {
			glog.Errorf("unable to get sig files: %v", err)
			return err
		}
		if lenSigFiles > 0 {
			glog.Infof("sig file is exist: %s", contents.Sha)
			// Check sha in database
			updatesf := database.SigFiles{}
			err = database.DBConnection.
				Where("owner = ? and repo = ? and path = ? and ref = ?", watchOwner, watchRepo, watchPath, watchRef).
				First(&updatesf).Error
			if err != nil {
				glog.Errorf("unable to get sig files: %v", err)
				return err
			}
			// write sha in waitingsha
			updatesf.WaitingSha = contents.Sha
			// init targetsha
			updatesf.TargetSha = ""
			err = database.DBConnection.Save(&updatesf).Error
			if err != nil {
				glog.Errorf("unable to get sig files: %v", err)
				return err
			}

		} else {
			glog.Infof("sig file is non-exist: %s", contents.Sha)
			// add sig file
			addsf := database.SigFiles{
				Owner:      watchOwner,
				Repo:       watchRepo,
				Path:       watchPath,
				Ref:        watchRef,
				WaitingSha: contents.Sha,
			}

			// create sig file
			err = database.DBConnection.Create(&addsf).Error
			if err != nil {
				glog.Errorf("unable to create sig files: %v", err)
				return err
			}
			glog.Infof("add sig file successfully: %s", contents.Sha)
		}
	}
	return nil
}

// watch database
func (handler *SigHandler) watch() {
	if len(handler.Config.WatchSigFiles) == 0 {
		return
	}

	for {
		watchDuration := handler.Config.WatchSigFileDuration
		for _, wf := range handler.Config.WatchSigFiles {
			// get params
			watchOwner := wf.WatchSigFileOwner
			watchRepo := wf.WatchSigFileRepo
			watchPath := wf.WatchSigFilePath
			watchRef := wf.WatchSigFileRef

			glog.Infof("begin to serve in sig. watchOwner: %s watchRepo: %s watchPath: %s watchRef: %s watchDuration: %d",
				watchOwner, watchRepo, watchPath, watchRef, watchDuration)

			// get sig file
			sf := database.SigFiles{}
			err := database.DBConnection.
				Where("owner = ? and repo = ? and path = ? and ref = ?", watchOwner, watchRepo, watchPath, watchRef).
				First(&sf).Error
			if err != nil {
				glog.Errorf("unable to get sig files: %v", err)
			} else {
				glog.Infof("init handler current sha: %v target sha: %v waiting sha: %v in sig",
					sf.CurrentSha, sf.TargetSha, sf.WaitingSha)
				if sf.TargetSha != "" {
					// skip when there is executing target sha
					glog.Infof("there is executing target sha: %v in sig", sf.TargetSha)
				} else {
					if sf.WaitingSha != "" && sf.CurrentSha != sf.WaitingSha {
						// waiting -> target
						sf.TargetSha = sf.WaitingSha
						err = database.DBConnection.Save(&sf).Error
						if err != nil {
							glog.Errorf("unable to save sig files: %v", err)
						} else {
							// define update sf
							updatesf := &database.SigFiles{}
							updatesf.ID = sf.ID

							// get file content from target sha
							glog.Infof("get target sha blob: %v", sf.TargetSha)
							localVarOptionals := &gitee.GetV5ReposOwnerRepoGitBlobsShaOpts{}
							localVarOptionals.AccessToken = optional.NewString(handler.Config.GiteeToken)
							blob, _, err := handler.GiteeClient.GitDataApi.GetV5ReposOwnerRepoGitBlobsSha(
								handler.Context, watchOwner, watchRepo, sf.TargetSha, localVarOptionals)
							if err != nil {
								glog.Errorf("unable to get blob: %v", err)
							} else {
								// base64 decode
								glog.Infof("decode target sha blob: %v", sf.TargetSha)
								decodeBytes, err := base64.StdEncoding.DecodeString(blob.Content)
								if err != nil {
									glog.Errorf("decode content with error: %v", err)
								} else {
									// unmarshal owners file
									glog.Infof("unmarshal target sha blob: %v", sf.TargetSha)
									var sy SigsYaml
									err = yaml.Unmarshal(decodeBytes, &sy)
									if err != nil {
										glog.Errorf("failed to unmarshal sigs: %v", err)
									} else {
										glog.Infof("get blob result: %v", sy)
										result := true
										// handle sigs
										err = handler.handleSigs(sy)
										if err != nil {
											glog.Errorf("failed to handle sig: %v", err)
											result = false
										}
										glog.Infof("running result: %v", result)
										if result {
											err = database.DBConnection.Model(updatesf).Update("CurrentSha", sf.TargetSha).Error
											if err != nil {
												glog.Errorf("unable to update current sha in sig: %v", err)
											}
										}
									}
								}
							}

							// at last update target sha
							err = database.DBConnection.Model(updatesf).Update("TargetSha", "").Error
							if err != nil {
								glog.Errorf("unable to update target sha in sig: %v", err)
							}
							glog.Info("update sha successfully in sig")
						}
					} else {
						glog.Infof("no waiting sha in sig: %v", sf.WaitingSha)
					}
				}
			}
		}

		// watch duration
		glog.Info("end to serve in sig")
		time.Sleep(time.Duration(watchDuration) * time.Second)
	}
}

// handleSigs handle sig
func (handler *SigHandler) handleSigs(sy SigsYaml) error {
	// handle sig repos
	for i := 0; i < len(sy.Sigs); i++ {
		// handle sig repos
		err := handler.handleSigRepos(sy.Sigs[i])
		if err != nil {
			glog.Errorf("failed to handle sig repos: %v", err)
			return err
		}
	}

	mapSigs := make(map[string]string)
	if len(sy.Sigs) > 0 {
		for _, s := range sy.Sigs {
			mapSigs[s.Name] = s.Name
		}
	}

	// get sigs from DB
	var srs []database.SigRecords
	err := database.DBConnection.Model(&database.SigRecords{}).Find(&srs).Error
	if err != nil {
		glog.Errorf("unable to get sig repos: %v", err)
		return err
	}
	mapSigsInDB := make(map[string]string)
	for _, sr := range srs {
		mapSigsInDB[sr.Name] = strconv.Itoa(int(sr.ID))
	}

	// remove
	err = handler.removeSigs(mapSigs, mapSigsInDB)
	if err != nil {
		glog.Errorf("unable to remove sig: %v", err)
	}

	// add
	err = handler.addSigs(mapSigs, mapSigsInDB)
	if err != nil {
		glog.Errorf("unable to add sig: %v", err)
	}

	return nil
}

// removeSigs
func (handler *SigHandler) removeSigs(mapSigs, mapSigsInDB map[string]string) error {
	listOfRemove := make([]string, 0)

	for k := range mapSigsInDB {
		if _, exists := mapSigs[k]; !exists {
			listOfRemove = append(listOfRemove, k)
		}
	}
	glog.Infof("list of remove sigs: %v", listOfRemove)

	if len(listOfRemove) > 0 {
		glog.Info("begin to remove sigs")
		for _, v := range listOfRemove {
			// remove from DB
			id, _ := strconv.Atoi(mapSigsInDB[v])
			sr := database.SigRecords{}
			sr.ID = uint(id)
			err := database.DBConnection.Delete(&sr).Error
			if err != nil {
				glog.Errorf("failed to remove sig in database: %v", err)
			}

			glog.Infof("begin to remove sig repos for %s", v)
			database.DBConnection.Where(database.SigRepositories{Name: v}).Delete(database.SigRepositories{})
			glog.Infof("end to remove sig repo for %s", v)
		}
		glog.Info("end to remove sigs")
	}

	return nil
}

// addSigs
func (handler *SigHandler) addSigs(mapSigs, mapSigsInDB map[string]string) error {
	listOfAdd := make([]string, 0)

	for k := range mapSigs {
		if _, exits := mapSigsInDB[k]; !exits {
			listOfAdd = append(listOfAdd, k)
		}
	}
	glog.Infof("list of add sigs: %v", listOfAdd)

	if len(listOfAdd) > 0 {
		glog.Info("begin to add sigs")
		for _, v := range listOfAdd {
			// add sig
			sr := database.SigRecords{
				Name: v,
			}
			err := database.DBConnection.Create(&sr).Error
			if err != nil {
				glog.Errorf("failed to add sigs in database: %v", err)
			}
		}
		glog.Info("end to add sigs")
	}

	return nil
}

// getSigsLength get sigs length
func (handler *SigHandler) getSigsLength(name string) (int, error) {
	// Check sigs file
	var lenSigs int
	err := database.DBConnection.Model(&database.SigRecords{}).
		Where("name = ?", name).
		Count(&lenSigs).Error
	if err != nil {
		glog.Errorf("unable to get sig length: %v", err)
	}
	return lenSigs, err
}

// handleSigRepos handle sig repos
func (handler *SigHandler) handleSigRepos(sig Sig) error {
	mapSigRepos := make(map[string]string)

	if len(sig.Repositories) > 0 {
		for _, sr := range sig.Repositories {
			mapSigRepos[sr] = sr
		}
	}

	// get sig repos from DB
	var srs []database.SigRepositories
	err := database.DBConnection.Model(&database.SigRepositories{}).
		Where("name = ?", sig.Name).Find(&srs).Error
	if err != nil {
		glog.Errorf("unable to get sig repos: %v", err)
		return err
	}
	mapSigReposInDB := make(map[string]string)
	for _, sr := range srs {
		mapSigReposInDB[sr.RepoName] = strconv.Itoa(int(sr.ID))
	}

	// remove
	err = handler.removeSigRepos(sig, mapSigRepos, mapSigReposInDB)
	if err != nil {
		glog.Errorf("unable to remove sig repos: %v", err)
	}

	// add
	err = handler.addSigRepos(sig, mapSigRepos, mapSigReposInDB)
	if err != nil {
		glog.Errorf("unable to add sig repos: %v", err)
	}

	return nil
}

// removeSigRepos
func (handler *SigHandler) removeSigRepos(sig Sig, mapSigRepos, mapSigReposInDB map[string]string) error {
	listOfRemove := make([]string, 0)

	for k := range mapSigReposInDB {
		if _, exists := mapSigRepos[k]; !exists {
			listOfRemove = append(listOfRemove, k)
		}
	}
	glog.Infof("list of remove sig repos: %v", listOfRemove)

	if len(listOfRemove) > 0 {
		glog.Infof("begin to remove sig repos for %s", sig.Name)
		for _, v := range listOfRemove {
			// remove from DB
			id, _ := strconv.Atoi(mapSigReposInDB[v])
			sr := database.SigRepositories{}
			sr.ID = uint(id)
			err := database.DBConnection.Delete(&sr).Error
			if err != nil {
				glog.Errorf("failed to remove sig repo in database: %v", err)
			}
		}
		glog.Infof("end to remove sig repo for %s", sig.Name)
	}

	return nil
}

// addSigRepos
func (handler *SigHandler) addSigRepos(sig Sig, mapSigRepos, mapSigReposInDB map[string]string) error {
	listOfAdd := make([]string, 0)

	for k := range mapSigRepos {
		if _, exits := mapSigReposInDB[k]; !exits {
			listOfAdd = append(listOfAdd, k)
		}
	}
	glog.Infof("list of add sig repos: %v", listOfAdd)

	if len(listOfAdd) > 0 {
		glog.Infof("begin to add sig repos for %s", sig.Name)
		for _, v := range listOfAdd {
			sr := database.SigRepositories{
				Name:     sig.Name,
				RepoName: v,
			}
			err := database.DBConnection.Create(&sr).Error
			if err != nil {
				glog.Errorf("failed to add sig repos in database: %v", err)
			}
		}
		glog.Infof("end to add sig repos for %s", sig.Name)
	}

	return nil
}

// get Sig name by Repo name
func (s *Server) getSigNameFromRepo(repoName string) (sigName string) {
	sigName = ""
	if len(repoName) == 0 || len(repoName) > 128 {
		glog.Errorf("Repo name is invalid.")
		return
	}
	// get sig repos from DB
	glog.Infof("Repo name is:%s .", repoName)
	var srs database.SigRepositories
	err := database.DBConnection.Model(&database.SigRepositories{}).
		Where("repo_name = ?", repoName).Find(&srs).Error
	if err != nil {
		glog.Errorf("unable to get sig repos: %v", err)
		return
	}
	sigName = srs.Name
	glog.Infof("end to add sig repos for %s", sigName)
	return
}

