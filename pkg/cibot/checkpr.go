package cibot

import (
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

var checkPrComment = "Cannot use \"/check-pr\", because this command is only used to detect open pull requests"

//CheckPr Check whether the pull request can be merged
func (s *Server) CheckPr(event *gitee.NoteEvent) (err error) {
	if *event.NoteableType == "PullRequest" && event.PullRequest.State == "open" {
		err := s.MergePullRequest(event)
		if err != nil {
			comment := err.Error()
			return s.addCommentToPullRequest(event.Repository.Namespace, event.Repository.Name,
				comment, event.PullRequest.Number)
		}
	} else {
		return s.addCommentToPullRequest(event.Repository.Namespace, event.Repository.Name,
			checkPrComment, event.PullRequest.Number)
	}
	return nil
}

func (s *Server) addCommentToPullRequest(owner, repo, comment string, number int32) error {
	body := gitee.PullRequestCommentPostParam{}
	body.AccessToken = s.Config.GiteeToken
	body.Body = comment
	_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
	if err != nil {
		glog.Errorf("unable to add comment in pull request: %v", err)
		return err
	}
	return nil
}
