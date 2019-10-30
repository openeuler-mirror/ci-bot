package cibot

import (
	"fmt"
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

const (
	issueAssignMessage       = `this issue is assigned to: ***%s***.`
	issueCanNotAssignMessage = `this issue can not be assigned to: ***%s***.
please try to assign to the repository collaborators.`
)

// Assign a collaborator for issue
func (s *Server) Assign(event *gitee.NoteEvent) error {
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
			glog.Infof("assign started. comment: %s owner: %s repo: %s issueNumber: %s issueAuthor: %s commentAuthor: %s",
				comment, owner, repo, issueNumber, issueAuthor, commentAuthor)

			// split the assignees and operation to be performed.
			substrings := strings.Split(comment, "@")
			// default is comment author
			assignee := commentAuthor
			if len(substrings) > 1 {
				assignee = strings.TrimSpace(substrings[1])
			}

			body := gitee.IssueUpdateParam{}
			body.Repo = repo
			body.Assignee = assignee
			body.AccessToken = s.Config.GiteeToken
			glog.Infof("invoke api to assign: %s", issueNumber)

			// patch assignee
			_, response, err := s.GiteeClient.IssuesApi.PatchV5ReposOwnerIssuesNumber(s.Context, owner, issueNumber, body)
			if err != nil {
				if response.StatusCode == 403 {
					glog.Infof("unable to assign with status code %d: %s", response.StatusCode, issueNumber)
					// add comment
					body := gitee.IssueCommentPostParam{}
					body.AccessToken = s.Config.GiteeToken
					body.Body = fmt.Sprintf(issueCanNotAssignMessage, assignee)
					_, _, err := s.GiteeClient.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(s.Context, owner, repo, issueNumber, body)
					if err != nil {
						glog.Errorf("unable to add comment in issue: %v", err)
					}
				} else {
					glog.Errorf("unable to assign: %s err: %v", issueNumber, err)
					return err
				}
			} else {
				glog.Infof("assign successfully: %v", issueNumber)
				// add comment
				body := gitee.IssueCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(issueAssignMessage, assignee)
				_, _, err := s.GiteeClient.IssuesApi.PostV5ReposOwnerRepoIssuesNumberComments(s.Context, owner, repo, issueNumber, body)
				if err != nil {
					glog.Errorf("unable to add comment in issue: %v", err)
				}
			}
		}
	}
	return nil
}
