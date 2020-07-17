package cibot

import (
	"fmt"
	"strings"

	"github.com/antihax/optional"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

const (
	lgtmSelfOwnMessage         = `***lgtm*** can not be added in your self-own pull request. :astonished: `
	lgtmAddedMessage           = `***lgtm*** is added in this pull request by: ***%s***. :wave: `
	lgtmRemovedMessage         = `***lgtm*** is removed in this pull request by: ***%s***. :flushed: `
	lgtmAddNoPermissionMessage = `***%s*** has no permission to add ***lgtm*** in this pull request. :astonished:
please contact to the collaborators in this repository.`
	lgtmRemoveNoPermissionMessage = `***%s*** has no permission to remove ***lgtm*** in this pull request. :astonished:
please contact to the collaborators in this repository.`
	lgtmRemovePullRequestChangeMessage = `new changes are detected. ***lgtm*** is removed in this pull request by: ***%s***. :flushed: `
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

			// check sigs
			r, err := canCommentPrIncludingSigDirectory(s, owner, repo, prNumber, commentAuthor)
			glog.Infof("sig owners check: can comment, r=%v, err=%v\n", r, err)
			if err != nil {
				glog.Errorf("unable to check sigs permission: %v", err)
				return err
			}

			// permission: admin, write, read, none
			if permission.Permission == "admin" || permission.Permission == "write" || isOwner || r == 1 {
				// add lgtm label
				addlabel := &gitee.NoteEvent{}
				addlabel.PullRequest = event.PullRequest
				addlabel.Repository = event.Repository
				addlabel.Comment = &gitee.NoteHook{}
				err = s.AddSpecifyLabelsInPulRequest(addlabel, []string{s.getLgtmLable(commentAuthor)}, true)
				if err != nil {
					return err
				}
				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(lgtmAddedMessage, commentAuthor) + fmt.Sprintf(LabelHiddenValue, event.PullRequest.Head.Sha)
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
				body.Body = fmt.Sprintf(lgtmAddNoPermissionMessage, commentAuthor)
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

func (s *Server) getLgtmLable(commenter string) string {
	if s.Config.LgtmCountsRequired > 1 {
		return fmt.Sprintf(LabelLgtmWithCommenter, strings.ToLower(commenter))
	}
	return LabelNameLgtm
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

				// check sigs
				r, err := canCommentPrIncludingSigDirectory(s, owner, repo, prNumber, commentAuthor)
				glog.Infof("sig owners check: can comment, r=%v, err=%v\n", r, err)
				if err != nil {
					glog.Errorf("unable to check sigs permission: %v", err)
					return err
				}

				// permission: admin, write, read, none
				if permission.Permission != "admin" && permission.Permission != "write" && !isOwner && r != 1 {
					glog.Info("no permission to remove lgtm")
					// add comment
					body := gitee.PullRequestCommentPostParam{}
					body.AccessToken = s.Config.GiteeToken
					body.Body = fmt.Sprintf(lgtmRemoveNoPermissionMessage, commentAuthor)
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
			}

			// remove lgtm label
			removelabel := &gitee.NoteEvent{}
			removelabel.PullRequest = event.PullRequest
			removelabel.Repository = event.Repository
			removelabel.Comment = &gitee.NoteHook{}
			mapOfRemoveLabels := map[string]string{}
			lgtmLable := s.getLgtmLable(commentAuthor)
			mapOfRemoveLabels[lgtmLable] = lgtmLable
			err := s.RemoveSpecifyLabelsInPulRequest(removelabel, mapOfRemoveLabels)
			if err != nil {
				return err
			}

			// add comment
			body := gitee.PullRequestCommentPostParam{}
			body.AccessToken = s.Config.GiteeToken
			body.Body = fmt.Sprintf(lgtmRemovedMessage, commentAuthor)
			number := event.PullRequest.Number
			_, _, err = s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
			if err != nil {
				glog.Errorf("unable to add comment in pull request: %v", err)
				return err
			}
		}
	}
	return nil
}

func (s *Server) collectExistingLgtmLabel(owner, repo string, number int32) (map[string]string, error) {
	labels := make(map[string]string)
	if s.Config.LgtmCountsRequired > 1 {
		lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{
			AccessToken: optional.NewString(s.Config.GiteeToken),
		}
		pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, lvos)
		if err != nil {
			glog.Errorf("unable to get pull request. err: %v", err)
			return nil, err
		}
		glog.Infof("list of pr labels: %v", pr.Labels)
		for _, label := range pr.Labels {
			if strings.HasPrefix(label.Name, fmt.Sprintf(LabelLgtmWithCommenter, "")) {
				labels[label.Name] = label.Name
			}
		}
	} else {
		labels[LabelNameLgtm] = LabelNameLgtm
	}
	return labels, nil
}

// CheckLgtmByPullRequestUpdate checks lgtm when received the pull request update event
func (s *Server) CheckLgtmByPullRequestUpdate(event *gitee.PullRequestEvent) error {
	owner := event.Repository.Namespace
	repo := event.Repository.Name
	prNumber := event.PullRequest.Number
	commentCount := event.PullRequest.Comments
	var perPage int32 = 20
	pageCount := commentCount / perPage
	if commentCount%perPage > 0 {
		pageCount++
	}
	if perPage == 0 {
		pageCount++
	}
	glog.Infof("pull request comment count: %v page count: %v per page: %v", commentCount, pageCount, perPage)

	// find comments from the last page
	lastlgtmSha := ""
	for page := pageCount; page > 0; page-- {
		localVarOptionals := &gitee.GetV5ReposOwnerRepoPullsNumberCommentsOpts{}
		localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)
		localVarOptionals.PerPage = optional.NewInt32(perPage)
		localVarOptionals.Page = optional.NewInt32(page)

		comments, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, prNumber, localVarOptionals)
		if err != nil {
			glog.Errorf("unable to get pull request comments. err: %v", err)
			return err
		}
		if len(comments) > 0 {
			// from the last comment
			for length := len(comments) - 1; length >= 0; length-- {
				comment := comments[length]
				m := RegBotAddLgtm.FindStringSubmatch(comment.Body)
				if /*comment.User.Login == BotName &&*/ m != nil && comment.UpdatedAt == comment.CreatedAt {
					lastlgtmSha = m[1]
					glog.Infof("pull request comment with lastlgtmSha: %v", comment)
					break
				}
			}
		}

		// get the last sha when lgtm
		if lastlgtmSha != "" {
			glog.Infof("lastlgtmSha: %v sha: %v", lastlgtmSha, event.PullRequest.Head.Sha)
			// if the sha is changed
			if lastlgtmSha != event.PullRequest.Head.Sha {
				// remove lgtm label
				removelabel := &gitee.NoteEvent{}
				removelabel.PullRequest = event.PullRequest
				removelabel.Repository = event.Repository
				removelabel.Comment = &gitee.NoteHook{}
				labels, err := s.collectExistingLgtmLabel(owner, repo, prNumber)
				if err != nil {
					glog.Errorf("unable to pick lgtm label in pr: %v", err)
					return err
				}
				err = s.RemoveSpecifyLabelsInPulRequest(removelabel, labels)
				if err != nil {
					return err
				}

				// add comment
				body := gitee.PullRequestCommentPostParam{}
				body.AccessToken = s.Config.GiteeToken
				body.Body = fmt.Sprintf(lgtmRemovePullRequestChangeMessage, s.Config.BotName)
				_, _, err = s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, prNumber, body)
				if err != nil {
					glog.Errorf("unable to add comment in pull request: %v", err)
					return err
				}
			}
			break
		}
	}

	return nil
}
