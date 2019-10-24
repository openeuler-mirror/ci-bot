package cibot

import (
	"encoding/base64"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var (
	DefaultOwnerFileName = "OWNERS"
)

type OwnersFile struct {
	Maintainers []string `yaml:"maintainers"`
}

// GetOwners gets owners from owners file in repository
func (s *Server) GetOwners(event *gitee.NoteEvent) []string {
	// get basic params
	owner := event.Repository.Namespace
	repo := event.Repository.Name
	branch := event.PullRequest.Base.Ref
	glog.Infof("get owners started. owner: %s repo: %s branch: %s", owner, repo, branch)

	localVarOptionals := &gitee.GetV5ReposOwnerRepoContentsPathOpts{}
	localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)
	localVarOptionals.Ref = optional.NewString(branch)
	// get contents
	contents, _, err := s.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
		s.Context, owner, repo, DefaultOwnerFileName, localVarOptionals)
	if err != nil {
		glog.Errorf("unable to get repository content by path: %v", err)
		return nil
	}

	if len(contents) > 0 {
		// base64 decode
		decodeBytes, err := base64.StdEncoding.DecodeString(contents[0].Content)
		if err != nil {
			glog.Errorf("decode content with error: %v", err)
			return nil
		}
		// unmarshal owners file
		var owners OwnersFile
		err = yaml.Unmarshal(decodeBytes, &owners)
		if err != nil {
			glog.Errorf("fail to unmarshal owners: %v", err)
		}

		// return owners
		if len(owners.Maintainers) > 0 {
			return owners.Maintainers
		}
	}

	return nil
}
