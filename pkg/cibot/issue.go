package cibot

import (
	"fmt"
	"strings"
	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
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

		owner := event.Repository.Namespace
		repo := event.Repository.Path
		number := event.Issue.Number
		// get sig name add a tag to describe the sig name of the repo.
		limitNotice := false
		sigName := s.getSigNameFromRepo(event.Repository.FullName)
		if len(sigName) > 0 {
			label := []string{fmt.Sprintf("sig/%s", sigName)}
			labelops := gitee.PullRequestLabelPostParam{s.Config.GiteeToken, label}
			_, _, error := s.GiteeClient.LabelsApi.PostV5ReposOwnerRepoIssuesNumberLabels(s.Context, owner, repo, number, labelops)
			if error != nil {
				glog.Errorf("unable to add label in issue: %v", error)
			}
			for _, limitSig := range s.Config.LimitMemberSigs {
				if sigName == limitSig {
					limitNotice = true
				}
			}
		}

		//get committor list:
		var ps []database.Privileges
		err := database.DBConnection.Model(&database.Privileges{}).
			Where("owner = ? and repo = ? and type = ?", owner, repo, PrivilegeDeveloper).Find(&ps).Error
		if err != nil {
			glog.Errorf("unable to get members: %v", err)
		}
		var committors []string
		if len(ps) > 0 {
			for _, p := range ps {
				committors = append(committors, fmt.Sprintf("***@%s***", p.User))
				if limitNotice && (len(committors) >= s.Config.LimitMemberCnt) {
					break
				}
			}
		}
		committor_list := strings.Join(committors, ", ")
		if len(committor_list) > 0 {
			sigPath := fmt.Sprintf(SigPath, sigName)
			proInfo := fmt.Sprintf(DisplayCommittors, sigName, sigPath)
			committor_list = proInfo + committor_list + "."
		}

		// add comment
		body := gitee.IssueCommentPostParam{}
		body.AccessToken = s.Config.GiteeToken
		body.Body = fmt.Sprintf(tipBotMessage, event.Sender.Login, s.Config.CommunityName, s.Config.CommandLink, committor_list)
		//Issue could exists without belonging to any repo.
		if event.Repository == nil {
			glog.Warningf("Issue is not created on repo, skip posting issue comment.")
			return
		}
		_, _, err = s.GiteeClient.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(s.Context, owner, repo, number, body)
		if err != nil {
			glog.Errorf("unable to add comment in issue: %v", err)
		}
	}
}
