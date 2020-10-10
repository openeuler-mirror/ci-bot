package cibot

import (
	"context"
	"encoding/base64"
	"errors"
	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"sync"
	"time"
)

//FrozenHandler Handling frozen branches
type FrozenHandler struct {
	Config      config.Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

type freezeFile struct {
	name string
	path string
	size string
	sha  string
}

//FrozenYaml  freeze configuration
type FrozenYaml struct {
	Release []FrozenBranchYaml `yaml:"release"`
}

//FrozenBranchYaml Branch freeze configuration
type FrozenBranchYaml struct {
	Branch string   `yaml:"branch"`
	Frozen bool     `yaml:"frozen"`
	Owner  []string `yaml:"owner"`
}

var
(
	frozenList []FrozenBranchYaml
	frozenFile freezeFile
	lock       sync.RWMutex
)

func (fh *FrozenHandler) Server() {
	err := fh.initFrozenFile()
	if err != nil {
		glog.Error(err)
	}
	go fh.watch()
}

func (fh *FrozenHandler) initFrozenFile() error {
	if fh.Config.WatchFrozenFile == (config.WatchFrozenFile{}) {
		return errors.New("Frozen configuration items are not initialized ")
	}
	fc, _, err := fh.getFrozenFileContent()
	if err != nil {
		return err
	}
	err = handleContent(fc)
	if err != nil {
		emptyFrozenList()
	}
	return err
}

func (fh *FrozenHandler) getFrozenFileContent() (content string, changed bool, err error) {
	localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
	localVarOptionals.AccessToken = optional.NewString(fh.Config.GiteeToken)
	localVarOptionals.Ref = optional.NewString(fh.Config.WatchFrozenFile.FrozenFileRef)
	contents, _, err := fh.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
		fh.Context, fh.Config.WatchFrozenFile.FrozenFileOwner, fh.Config.WatchFrozenFile.FrozenFileRepo,
		fh.Config.WatchFrozenFile.FrozenFilePath, localVarOptionals)
	if err != nil {
		return "", changed, err
	}
	if frozenFile.sha != contents.Sha {
		frozenFile.path = contents.Path
		frozenFile.name = contents.Name
		frozenFile.size = contents.Size
		frozenFile.sha = contents.Sha
		changed = true
	}
	return contents.Content, changed, nil
}

func (fh *FrozenHandler) watch() {
	if fh.Config.WatchFrozenFile == (config.WatchFrozenFile{}) {
		return
	}
	watchDuration := fh.Config.WatchFrozenDuration
	for {
		fileContent, changed, err := fh.getFrozenFileContent()
		if err != nil {
			emptyFrozenList()
			glog.Error(err)
		} else {
			if changed {
				err := handleContent(fileContent)
				if err != nil {
					glog.Error(err)
					emptyFrozenList()
				}
			}
		}
		time.Sleep(time.Duration(watchDuration) * time.Second)
	}
}

func handleContent(content string) error {
	if content == "" {
		return errors.New("The parsed content cannot be empty ")
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}
	var fz FrozenYaml
	err = yaml.Unmarshal(decodeBytes, &fz)
	if err != nil {
		return err
	}
	if len(fz.Release) == 0 {
		glog.Info("unmarshal frozen branch is empty")
		emptyFrozenList()
	} else {
		extractFrozenBranch(fz.Release)
	}
	return nil
}

func extractFrozenBranch(release []FrozenBranchYaml) {
	var fs []FrozenBranchYaml
	for _, v := range release {
		if v.Frozen {
			fs = append(fs, v)
		}
	}
	if len(fs) > 0 {
		writeFrozenList(fs)
	} else {
		emptyFrozenList()
	}
}
func emptyFrozenList() {
	lock.Lock()
	if len(frozenList) > 0 {
		frozenList = frozenList[:0]
	}
	lock.Unlock()
}
func writeFrozenList(fs []FrozenBranchYaml) {
	lock.Lock()
	defer lock.Unlock()
	if len(frozenList) > 0 {
		frozenList = frozenList[:0]
	}
	frozenList = append(frozenList, fs...)
}

//IsBranchFrozen Check if the branch is frozen
func IsBranchFrozen(branch string) (owner []string, isFrozen bool) {
	lock.RLock()
	defer lock.RUnlock()
	if len(frozenList) == 0 {
		return nil, isFrozen
	}
	for _, v := range frozenList {
		if v.Branch == branch {
			isFrozen = true
			owner = append(owner, v.Owner...)
			break
		}
	}
	return
}
