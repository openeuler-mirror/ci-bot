package cibot

import (
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

// HandlePullRequestEvent handles pull request event
func (s *Server) HandlePullRequestEvent(event *gitee.PullRequestEvent) {
	if event == nil {
		return
	}

	glog.Infof("pull request sender: %v", event.Sender.Login)

	// handle events
	switch *event.Action {
	case "create":
		glog.Info("received a pull request create event")
		err := s.CheckCLAByPullRequestEvent(event)
		if err != nil {
			glog.Errorf("failed to check cla by pull request event: %v", err)
		}
	}
}

// MergePullRequest with lgtm and approved label
func (s *Server) MergePullRequest(event *gitee.NoteEvent) error {
	// get basic params
	owner := event.Repository.Owner.Login
	repo := event.Repository.Name
	prNumber := event.PullRequest.Number
	glog.Infof("merge pull request started. owner: %s repo: %s number: %d", owner, repo, prNumber)

	// list labels in current pull request
	pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, prNumber, nil)
	if err != nil {
		glog.Errorf("unable to get pull request. err: %v", err)
		return err
	}
	listofPrLabels := pr.Labels
	glog.Infof("List of pr labels: %v", listofPrLabels)

	// check if it has both lgtm and approved label
	hasApproved := false
	hasLgtm := false
	for _, l := range listofPrLabels {
		if l.Name == LabelNameLgtm {
			hasLgtm = true
		} else if l.Name == LabelNameApproved {
			hasApproved = true
		}
	}
	glog.Infof("Pr labels have approved: %t lgtm: %t", hasApproved, hasLgtm)

	// ready to merge
	if hasApproved && hasLgtm {
		// current pr can be merged
		if event.PullRequest.Mergeable {
			// merge pr
			body := gitee.PullRequestMergePutParam{}
			body.AccessToken = s.Config.GiteeToken
			_, err = s.GiteeClient.PullRequestsApi.PutV5ReposOwnerRepoPullsNumberMerge(s.Context, owner, repo, prNumber, body)
			if err != nil {
				glog.Errorf("unable to merge pull request. err: %v", err)
				return err
			}
		}
	}

	return nil
}
