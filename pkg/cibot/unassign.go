package cibot

import (
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
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
				body.Assignee = ""
				body.AccessToken = s.Config.GiteeToken
				glog.Infof("invoke api to unassign: %s", issueNumber)

				// patch assignee
				_, _, err := s.GiteeClient.IssuesApi.PatchV5ReposOwnerIssuesNumber(s.Context, owner, issueNumber, body)
				if err != nil {
					glog.Errorf("unable to unassign: %s err: %v", issueNumber, err)
					return err
				}
				glog.Infof("unassign successfully: %v", issueNumber)
			} else {
				glog.Infof("no need to unassign: %v", issueNumber)
			}
		}
	}
	return nil
}
