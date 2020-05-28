package cibot

import (
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
)

func (s *Server) hasSquashCommitLabel(labels []gitee.Label) bool {
	for _, l := range labels {
		if s.Config.SquashCommitLabel == l.Name {
			return true
		}
	}
	return false
}

func (s *Server) ValidateCommits(event *gitee.PullRequestEvent) error {
	if s.Config.CommitsThreshold == 0 || len(s.Config.SquashCommitLabel) == 0 {
		glog.Info("'commitsthreshold' or 'squashcommitlabel' is not configured, skip validating pr commits")
		return nil
	}
	// Get Pull Request Commits detail
	owner := event.Repository.Namespace
	repo := event.Repository.Name
	prNumber := event.PullRequest.Number
	commitPullRequestOpts := &gitee.GetV5ReposOwnerRepoPullsNumberCommitsOpts{}
	commitPullRequestOpts.AccessToken = optional.NewString(s.Config.GiteeToken)
	commits, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberCommits(
		s.Context, owner, repo, prNumber, commitPullRequestOpts)
	if err != nil {
		glog.Errorf("failed to get pull request commits detail : %v", err)
		return err
	}
	commitsExceeded := len(commits) > s.Config.CommitsThreshold

	// Get Pull Request Label details
	prOpts := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{}
	prOpts.AccessToken = optional.NewString(s.Config.GiteeToken)
	pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, prNumber, prOpts)
	if err != nil {
		glog.Errorf("unable to get pull request. err: %v", err)
		return err
	}

	labelParam := &gitee.NoteEvent{}
	labelParam.PullRequest = event.PullRequest
	labelParam.Repository = event.Repository
	labelParam.Comment = &gitee.NoteHook{}
	if commitsExceeded && !s.hasSquashCommitLabel(pr.Labels) {
		//add squash commits label if needed
		err = s.AddSpecifyLabelsInPulRequest(labelParam, []string{s.Config.SquashCommitLabel}, true)
		if err != nil {
			return err
		}
	} else if !commitsExceeded && s.hasSquashCommitLabel(pr.Labels) {
		//remove label if needed
		mapOfRemoveLabels := map[string]string{}
		mapOfRemoveLabels[s.Config.SquashCommitLabel] = s.Config.SquashCommitLabel
		err = s.RemoveSpecifyLabelsInPulRequest(labelParam, mapOfRemoveLabels)
		if err != nil {
			return err
		}
	}
	return nil
}
