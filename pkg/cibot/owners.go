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

// CheckIsOwner checks the author is owner in repository
func (s *Server) CheckIsOwner(event *gitee.NoteEvent, author string) bool {
	isOwner := false
	// get owners
	owners := s.GetOwners(event)
	if owners != nil {
		glog.Infof("check isowner started. owners: %v", owners)
		if len(owners) > 0 {
			for _, owner := range owners {
				if owner == author {
					isOwner = true
					break
				}
			}
		}
	}
	glog.Infof("check %s isowner result: %t", author, isOwner)
	return isOwner
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

	return DecodeOwners(contents.Content)
}

func DecodeOwners(content string) []string {
	// base64 decode
	decodeBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		glog.Errorf("decode content with error: %v", err)
		return nil
	}
	// unmarshal owners file
	var owners OwnersFile
	err = yaml.Unmarshal(decodeBytes, &owners)
	if err != nil {
		glog.Errorf("fail to unmarshal owners: %v", err)
		return nil
	}

	// return owners
	if len(owners.Maintainers) > 0 {
		return owners.Maintainers
	}

	return nil
}
