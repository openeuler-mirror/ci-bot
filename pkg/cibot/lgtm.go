package cibot

import (
	"github.com/antihax/optional"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

const (
	lgtmSelfOwnMessage = `Thanks for your comment, but you can not add **lgtm** on your self-own pull request.:astonished: `
)

// AddLgtm adds lgtm label
func (s *Server) AddLgtm(event *gitee.NoteEvent) error {
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
			glog.Infof("add lgtm started. comment: %s prAuthor: %s commentAuthor: %s owner: %s repo: %s number: %d",
				comment, prAuthor, commentAuthor, owner, repo, prNumber)

			// can not lgtm on self-own pr
			if prAuthor == commentAuthor {
				glog.Info("can not lgtm on self-own pr")
				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = lgtmSelfOwnMessage
				owner := event.Repository.Namespace
				repo := event.Repository.Name
				number := event.PullRequest.Number
				_, _, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
				if err != nil {
					glog.Errorf("unable to add comment in pull request: %v", err)
					return err
				}
				return nil
			}

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
				// add lgtm label
				addlabel := &gitee.NoteEvent{}
				addlabel.PullRequest = event.PullRequest
				addlabel.Repository = event.Repository
				addlabel.Comment = &gitee.Note{}
				mapOfAddLabels := map[string]string{}
				mapOfAddLabels[LabelNameLgtm] = LabelNameLgtm
				err = s.AddSpecifyLabelsInPulRequest(addlabel, mapOfAddLabels)
				if err != nil {
					return err
				}
				// try to merge pr
				err = s.MergePullRequest(event)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// RemoveLgtm removes lgtm label
func (s *Server) RemoveLgtm(event *gitee.NoteEvent) error {
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
			glog.Infof("remove lgtm started. comment: %s prAuthor: %s commentAuthor: %s owner: %s repo: %s number: %d",
				comment, prAuthor, commentAuthor, owner, repo, prNumber)

			// can cancel lgtm on self-own pr
			if prAuthor != commentAuthor {
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
				if permission.Permission != "admin" && permission.Permission != "write" && !isOwner {
					glog.Info("no permission to remove lgtm")
					return nil
				}
			}

			// remove lgtm label
			removelabel := &gitee.NoteEvent{}
			removelabel.PullRequest = event.PullRequest
			removelabel.Repository = event.Repository
			removelabel.Comment = &gitee.Note{}
			mapOfRemoveLabels := map[string]string{}
			mapOfRemoveLabels[LabelNameLgtm] = LabelNameLgtm
			err := s.RemoveSpecifyLabelsInPulRequest(removelabel, mapOfRemoveLabels)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
