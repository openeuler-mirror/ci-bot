package cibot

import (
	"fmt"
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
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
	case "open":
		glog.Info("received a pull request open event")

		// add comment
		body := gitee.PullRequestCommentPostParam{}
		body.AccessToken = s.Config.GiteeToken
		body.Body = fmt.Sprintf(tipBotMessage, event.Sender.Login, s.Config.CommunityName, s.Config.CommunityName,
			s.Config.BotName, s.Config.CommandLink)
		owner := event.Repository.Namespace
		repo := event.Repository.Name
		number := event.PullRequest.Number
		_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
		if err != nil {
			glog.Errorf("unable to add comment in pull request: %v", err)
		}

		err = s.CheckCLAByPullRequestEvent(event)
		if err != nil {
			glog.Errorf("failed to check cla by pull request event: %v", err)
		}
	case "update":
		glog.Info("received a pull request update event")

		// get pr info
		owner := event.Repository.Namespace
		repo := event.Repository.Name
		number := event.PullRequest.Number
		lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{}
		lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
		pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, lvos)
		if err != nil {
			glog.Errorf("unable to get pull request. err: %v", err)
			return
		}
		listofPrLabels := pr.Labels
		glog.Infof("List of pr labels: %v", listofPrLabels)

		// remove lgtm if changes happen
		if s.hasLgtmLabel(pr.Labels) {
			err = s.CheckLgtmByPullRequestUpdate(event)
			if err != nil {
				glog.Errorf("check lgtm by pull request update. err: %v", err)
				return
			}
		}
	}
}

// RemoveAssigneesInPullRequest remove assignees in pull request
func (s *Server) RemoveAssigneesInPullRequest(event *gitee.NoteEvent) error {
	if event != nil {
		if event.PullRequest != nil {
			assignees := event.PullRequest.Assignees
			glog.Infof("remove assignees: %v", assignees)
			if len(assignees) > 0 {
				var strAssignees string
				for _, assignee := range assignees {
					strAssignees += assignee.Login + ","
				}
				strAssignees = strings.TrimRight(strAssignees, ",")
				glog.Infof("remove assignees str: %s", strAssignees)

				// get basic params
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				prNumber := event.PullRequest.Number
				localVarOptionals := &gitee.DeleteV5ReposOwnerRepoPullsNumberAssigneesOpts{}
				localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)

				// invoke api
				_, _, err := s.GiteeClient.PullRequestsApi.DeleteV5ReposOwnerRepoPullsNumberAssignees(s.Context, owner, repo, prNumber, strAssignees, localVarOptionals)
				if err != nil {
					glog.Errorf("unable to remove assignees in pull request. err: %v", err)
					return err
				}
				glog.Infof("remove assignees successfully: %s", strAssignees)
			}
		}
	}
	return nil
}

// RemoveTestersInPullRequest remove testers in pull request
func (s *Server) RemoveTestersInPullRequest(event *gitee.NoteEvent) error {
	if event != nil {
		if event.PullRequest != nil {
			testers := event.PullRequest.Testers
			glog.Infof("remove testers: %v", testers)
			if len(testers) > 0 {
				var strTesters string
				for _, tester := range testers {
					strTesters += tester.Login + ","
				}
				strTesters = strings.TrimRight(strTesters, ",")
				glog.Infof("remove testers str: %s", strTesters)

				// get basic params
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				prNumber := event.PullRequest.Number
				localVarOptionals := &gitee.DeleteV5ReposOwnerRepoPullsNumberTestersOpts{}
				localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)

				// invoke api
				_, _, err := s.GiteeClient.PullRequestsApi.DeleteV5ReposOwnerRepoPullsNumberTesters(s.Context, owner, repo, prNumber, strTesters, localVarOptionals)
				if err != nil {
					glog.Errorf("unable to remove testers in pull request. err: %v", err)
					return err
				}
				glog.Infof("remove testers successfully: %s", strTesters)
			}
		}
	}
	return nil
}

func (s *Server) hasLgtmLabel(labels []gitee.Label) bool {
	for _, l := range labels {
		if strings.HasPrefix(l.Name, fmt.Sprintf(LabelLgtmWithCommenter, "")) || l.Name == LabelNameLgtm {
			return true
		}
	}
	return false
}

func (s *Server) legalForMerge(labels []gitee.Label) bool {
	aproveLabel := 0
	lgtmLabel := 0
	lgtmPrefix := ""
	leastLgtm := 0
	if s.Config.LgtmCountsRequired > 1 {
		leastLgtm = s.Config.LgtmCountsRequired
		lgtmPrefix =fmt.Sprintf(LabelLgtmWithCommenter, "")
	} else {
		leastLgtm = 1
		lgtmPrefix = LabelNameLgtm
	}
	for _, l := range labels {
		if strings.HasPrefix(l.Name, lgtmPrefix) {
			lgtmLabel++
		} else if l.Name == LabelNameApproved {
			aproveLabel++
		}
	}
	glog.Infof("Pr labels have approved: %d lgtm: %d, required (%d)", aproveLabel, lgtmLabel, leastLgtm)
	return aproveLabel == 1 && lgtmLabel >= leastLgtm
}

// MergePullRequest with lgtm and approved label
func (s *Server) MergePullRequest(event *gitee.NoteEvent) error {
	// get basic params
	owner := event.Repository.Namespace
	repo := event.Repository.Name
	prNumber := event.PullRequest.Number
	glog.Infof("merge pull request started. owner: %s repo: %s number: %d", owner, repo, prNumber)

	// list labels in current pull request
	lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{}
	lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
	pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, prNumber, lvos)
	if err != nil {
		glog.Errorf("unable to get pull request. err: %v", err)
		return err
	}
	listofPrLabels := pr.Labels
	glog.Infof("List of pr labels: %v", listofPrLabels)

	// ready to merge
	if s.legalForMerge(listofPrLabels) {
		// current pr can be merged
		if event.PullRequest.Mergeable {
			// remove assignees
			err = s.RemoveAssigneesInPullRequest(event)
			if err != nil {
				glog.Errorf("unable to remove assignees. err: %v", err)
				return err
			}
			// remove testers
			err = s.RemoveTestersInPullRequest(event)
			if err != nil {
				glog.Errorf("unable to remove testers. err: %v", err)
				return err
			}
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
