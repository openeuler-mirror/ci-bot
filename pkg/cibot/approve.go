package cibot

import (
	"fmt"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
)

const (
	approvedAddedMessage           = `***approved*** is added in this pull request by: ***%s***. :wave: `
	approvedRemovedMessage         = `***approved*** is removed in this pull request by: ***%s***. :flushed: `
	approvedAddNoPermissionMessage = `***%s*** has no permission to add ***approved*** in this pull request. :astonished:
please contact to the collaborators in this repository.`
	approvedRemoveNoPermissionMessage = `***%s*** has no permission to remove ***approved*** in this pull request. :astonished:
please contact to the collaborators in this repository.`
)

// AddApprove adds approved label
func (s *Server) AddApprove(event *gitee.NoteEvent) error {
	// handle PullRequest
	if *event.NoteableType == "PullRequest" {
		// handle open
		if event.PullRequest.State == "open" {
			// get basic params
			comment := event.Comment.Body
			owner := event.Repository.Namespace
			repo := event.Repository.Name
			prAuthor := event.PullRequest.User.Login
			prNumber := event.PullRequest.Number
			commentAuthor := event.Comment.User.Login
			glog.Infof("add approve started. comment: %s prAuthor: %s commentAuthor: %s owner: %s repo: %s number: %d",
				comment, prAuthor, commentAuthor, owner, repo, prNumber)

			// check if current author has write permission
			localVarOptionals := &gitee.GetV5ReposOwnerRepoCollaboratorsUsernamePermissionOpts{}
			localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)
			// get permission
			permission, _, err := s.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoCollaboratorsUsernamePermission(
				s.Context, owner, repo, commentAuthor, localVarOptionals)
			if err != nil {
				glog.Errorf("unable to get comment author permission: %v", err)
				return err
			}
			// check author is owner
			isOwner := s.CheckIsOwner(event, commentAuthor)
			// permission: admin, write, read, none
			if permission.Permission == "admin" || permission.Permission == "write" || isOwner {
				// add approved label
				addlabel := &gitee.NoteEvent{}
				addlabel.PullRequest = event.PullRequest
				addlabel.Repository = event.Repository
				addlabel.Comment = &gitee.Note{}
				mapOfAddLabels := map[string]string{}
				mapOfAddLabels[LabelNameApproved] = LabelNameApproved
				err = s.AddSpecifyLabelsInPulRequest(addlabel, mapOfAddLabels)
				if err != nil {
					return err
				}
				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(approvedAddedMessage, commentAuthor)
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				number := event.PullRequest.Number
				_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
				if err != nil {
					glog.Errorf("unable to add comment in pull request: %v", err)
					return err
				}
				// try to merge pr
				err = s.MergePullRequest(event)
				if err != nil {
					return err
				}
			} else {
				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(approvedAddNoPermissionMessage, commentAuthor)
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				number := event.PullRequest.Number
				_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
				if err != nil {
					glog.Errorf("unable to add comment in pull request: %v", err)
					return err
				}
			}
		}
	}
	return nil
}

// RemoveApprove removes approved label
func (s *Server) RemoveApprove(event *gitee.NoteEvent) error {
	// handle PullRequest
	if *event.NoteableType == "PullRequest" {
		// handle open
		if event.PullRequest.State == "open" {
			// get basic params
			comment := event.Comment.Body
			owner := event.Repository.Namespace
			repo := event.Repository.Name
			prAuthor := event.PullRequest.User.Login
			prNumber := event.PullRequest.Number
			commentAuthor := event.Comment.User.Login
			glog.Infof("remove approve started. comment: %s prAuthor: %s commentAuthor: %s owner: %s repo: %s number: %d",
				comment, prAuthor, commentAuthor, owner, repo, prNumber)

			// check if current author has write permission
			localVarOptionals := &gitee.GetV5ReposOwnerRepoCollaboratorsUsernamePermissionOpts{}
			localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)
			// get permission
			permission, _, err := s.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoCollaboratorsUsernamePermission(
				s.Context, owner, repo, commentAuthor, localVarOptionals)
			if err != nil {
				glog.Errorf("unable to get comment author permission: %v", err)
				return err
			}
			// check author is owner
			isOwner := s.CheckIsOwner(event, commentAuthor)
			// permission: admin, write, read, none
			if permission.Permission == "admin" || permission.Permission == "write" || isOwner {
				// remove approved label
				removelabel := &gitee.NoteEvent{}
				removelabel.PullRequest = event.PullRequest
				removelabel.Repository = event.Repository
				removelabel.Comment = &gitee.Note{}
				mapOfRemoveLabels := map[string]string{}
				mapOfRemoveLabels[LabelNameApproved] = LabelNameApproved
				err = s.RemoveSpecifyLabelsInPulRequest(removelabel, mapOfRemoveLabels)
				if err != nil {
					return err
				}
				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(approvedRemovedMessage, commentAuthor)
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				number := event.PullRequest.Number
				_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
				if err != nil {
					glog.Errorf("unable to add comment in pull request: %v", err)
					return err
				}
			} else {
				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(approvedRemoveNoPermissionMessage, commentAuthor)
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				number := event.PullRequest.Number
				_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
				if err != nil {
					glog.Errorf("unable to add comment in pull request: %v", err)
					return err
				}
			}
		}
	}
	return nil
}
