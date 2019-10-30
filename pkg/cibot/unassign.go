package cibot

import (
	"fmt"
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

const (
	issueUnAssignMessage       = `***@%s*** is unassigned from this issue.`
	issueCanNotUnAssignMessage = `***@%s*** can not be unassigned from this issue.
please try to unassign the assignee from this issue.`
)

// UnAssign a collaborator for issue
func (s *Server) UnAssign(event *gitee.NoteEvent) error {
	if *event.NoteableType == "Issue" {
		// handle open
		if event.Issue.State == "open" {
			// get basic informations
			comment := event.Comment.Body
			owner := event.Repository.Namespace
			repo := event.Repository.Name
			issueNumber := event.Issue.Number
			issueAuthor := event.Issue.User.Login
			commentAuthor := event.Comment.User.Login
			glog.Infof("unassign started. comment: %s owner: %s repo: %s issueNumber: %s issueAuthor: %s commentAuthor: %s",
				comment, owner, repo, issueNumber, issueAuthor, commentAuthor)

			// split the assignees and operation to be performed.
			substrings := strings.Split(comment, "@")
			// default is comment author
			unassignee := commentAuthor
			if len(substrings) > 1 {
				unassignee = strings.TrimSpace(substrings[1])
			}

			// if issue assignee is the same with unassignee
			issueAssignee := ""
			if event.Issue.Assignee != nil {
				issueAssignee = event.Issue.Assignee.Login
			}
			if issueAssignee == unassignee {
				body := gitee.IssueUpdateParam{}
				body.Repo = repo
				// todo: the assignee can no be empty, this is a bug from gitee.
				// we need to change this once gitee fixes this bug.
				body.Assignee = ""
				body.AccessToken = s.Config.GiteeToken
				// build label string
				var strLabel string
				for _, l := range event.Issue.Labels {
					strLabel += l.Name + ","
				}
				strLabel = strings.TrimRight(strLabel, ",")
				if strLabel == "" {
					strLabel = ","
				}
				body.Labels = strLabel
				glog.Infof("invoke api to unassign: %s", issueNumber)

				// patch assignee
				_, _, err := s.GiteeClient.IssuesApi.PatchV5ReposOwnerIssuesNumber(s.Context, owner, issueNumber, body)
				if err != nil {
					glog.Errorf("unable to unassign: %s err: %v", issueNumber, err)
					return err
				}
				glog.Infof("unassign successfully: %v", issueNumber)

				// add comment
				bodyComment := gitee.IssueCommentPostParam{}
				bodyComment.AccessToken = s.Config.GiteeToken
				bodyComment.Body = fmt.Sprintf(issueUnAssignMessage, unassignee)
				_, _, err = s.GiteeClient.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(s.Context, owner, repo, issueNumber, bodyComment)
				if err != nil {
					glog.Errorf("unable to add comment in issue: %v", err)
				}
			} else {
				glog.Infof("can not unassign: %v", issueNumber)
				// add comment
				body := gitee.IssueCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(issueCanNotUnAssignMessage, unassignee)
				_, _, err := s.GiteeClient.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(s.Context, owner, repo, issueNumber, body)
				if err != nil {
					glog.Errorf("unable to add comment in issue: %v", err)
				}
			}
		}
	}
	return nil
}
