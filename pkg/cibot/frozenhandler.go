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
	"strings"
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
	if len(fh.Config.WatchFrozenFile) == 0 {
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

func (fh *FrozenHandler) getFrozenFileContent() (content []string, changed bool, err error) {
	localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
	localVarOptionals.AccessToken = optional.NewString(fh.Config.GiteeToken)
	for _,v :=range  fh.Config.WatchFrozenFile {
		localVarOptionals.Ref = optional.NewString(v.FrozenFileRef)
		contents, _, err := fh.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
			fh.Context, v.FrozenFileOwner, v.FrozenFileRepo,
			v.FrozenFilePath, localVarOptionals)
		if err != nil {
			glog.Error(err)
			continue
		}
		if !strings.Contains(frozenFile.sha, contents.Sha) {
			frozenFile.sha += ";"+contents.Sha
			changed = true
		}
		content = append(content, contents.Content)

	}
	if len(content) == 0 {
	   return content,changed,errors.New("Freeze information not obtained. ")
	}
	return content, changed, nil
}

func (fh *FrozenHandler) watch() {
	if len(fh.Config.WatchFrozenFile) == 0 {
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

func handleContent(content []string) error {
	if len(content) == 0 {
		return errors.New("The parsed content cannot be empty ")
	}
	var fzList []FrozenBranchYaml
	for _,v := range content {
		decodeBytes, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			glog.Error(err)
			continue
		}
		var fz FrozenYaml
		err = yaml.Unmarshal(decodeBytes, &fz)
		if err != nil {
			glog.Error(err)
			continue
		}
		if len(fz.Release) == 0 {
			glog.Info("unmarshal frozen branch is empty")
		}else {
			fzList = append(fzList, fz.Release...)
		}
	}

	if len(fzList) == 0 {
		return errors.New("all of frozen file unmarshal is empty or fail. ")
	}
	extractFrozenBranch(fzList)
	return nil
}

func extractFrozenBranch(release []FrozenBranchYaml) {
	emptyFrozenList()
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
		frozenFile = freezeFile{}
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
