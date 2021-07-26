package cibot

import (
	"fmt"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

// HandleIssueEvent handles issue event
func (s *Server) HandleIssueEvent(event *gitee.IssueEvent) {
	if event == nil {
		return
	}

	// handle events
	switch *event.Action {
	case "open":
		glog.Info("received a issue open event")

		// add comment
		body := gitee.IssueCommentPostParam{}
		body.AccessToken = s.Config.GiteeToken
		body.Body = fmt.Sprintf(tipBotMessage, event.Sender.Login, s.Config.CommunityName, s.Config.CommunityName,
			s.Config.BotName, s.Config.CommandLink)
		//Issue could exists without belonging to any repo.
		if event.Repository == nil {
			glog.Warningf("Issue is not created on repo, skip posting issue comment.")
			return
		}
		owner := event.Repository.Namespace
		repo := event.Repository.Path
		number := event.Issue.Number
		_, _, err := s.GiteeClient.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(s.Context, owner, repo, number, body)
		if err != nil {
			glog.Errorf("unable to add comment in issue: %v", err)
		}

		sigName := s.getSigNameFromRepo(event.Repository.FullName)
		if len(sigName) > 0 {
			label := []string{fmt.Sprintf("sig/%s", sigName[0:16])}
			labelops := gitee.PullRequestLabelPostParam{s.Config.GiteeToken, label}
			_, _, err = s.GiteeClient.LabelsApi.PostV5ReposOwnerRepoIssuesNumberLabels(s.Context, owner, repo, number, labelops)
			if err != nil {
				glog.Errorf("unable to add label in issue: %v", err)
			}
		}
	}
}
