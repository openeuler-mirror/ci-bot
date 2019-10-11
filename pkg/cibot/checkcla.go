package cibot

import (
	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

const (
	claNotFoundMessage = `Thanks for your pull request.
**Before we can look at your pull request, you'll need to sign a Contributor License Agreement (CLA).**
**Please follow instructions at <https://openeuler.org/en/cla.html> to sign the CLA.**
It may take a couple minutes for the CLA signature to be fully registered;
after that, please reply here with a new comment **/check-cla** and we'll verify.
- If you've already signed a CLA, it's possible we don't have your Gitee username or you're using a different email address.
  Check your existing CLA data and verify that your email at <https://gitee.com/profile/emails>.
- If you have done the above and are still having issues with the CLA being reported as unsigned,
  send a message to the backup e-mail support address at: contact@openeuler.org
`
)

// CheckCLAByNoteEvent check cla by NoteEvent
func (s *Server) CheckCLAByNoteEvent(event *gitee.NoteEvent) error {
	if *event.NoteableType == "PullRequest" {
		// PullRequest
		email := event.PullRequest.User.Email
		// check the email from sender
		var lenEmail int
		err := database.DBConnection.Model(&database.CLADetails{}).
			Where("email = ?", email).Count(&lenEmail).Error
		if err != nil {
			glog.Errorf("failed to check user email: %v", err)
			return err
		}
		if lenEmail > 0 {
			// add label openeuler-cla/yes
			addlabel := &gitee.NoteEvent{}
			addlabel.PullRequest = event.PullRequest
			addlabel.Repository = event.Repository
			addlabel.Comment = &gitee.Note{}
			addlabel.Comment.Body = AddClaYes
			err = s.AddLabelInPulRequest(addlabel)
			if err != nil {
				return err
			}

			// remove label openeuler-cla/no
			removelabel := &gitee.NoteEvent{}
			removelabel.PullRequest = event.PullRequest
			removelabel.Repository = event.Repository
			removelabel.Comment = &gitee.Note{}
			removelabel.Comment.Body = RemoveClaNo
			err = s.RemoveLabelInPullRequest(removelabel)
			if err != nil {
				return err
			}
		} else {
			// add label openeuler-cla/no
			addlabel := &gitee.NoteEvent{}
			addlabel.PullRequest = event.PullRequest
			addlabel.Repository = event.Repository
			addlabel.Comment = &gitee.Note{}
			addlabel.Comment.Body = AddClaNo
			err = s.AddLabelInPulRequest(addlabel)
			if err != nil {
				return err
			}

			// remove label openeuler-cla/yes
			removelabel := &gitee.NoteEvent{}
			removelabel.PullRequest = event.PullRequest
			removelabel.Repository = event.Repository
			removelabel.Comment = &gitee.Note{}
			removelabel.Comment.Body = RemoveClaYes
			err = s.RemoveLabelInPullRequest(removelabel)
			if err != nil {
				return err
			}

			// add comment
			body := gitee.PullRequestCommentPostParam{}
			body.AccessToken = s.Config.GiteeToken
			body.Body = claNotFoundMessage
			owner := event.Repository.Owner.Login
			repo := event.Repository.Name
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

// CheckCLAByPullRequestEvent check cla by PullRequestEvent
func (s *Server) CheckCLAByPullRequestEvent(event *gitee.PullRequestEvent) error {
	// check the email from sender
	var lenEmail int
	err := database.DBConnection.Model(&database.CLADetails{}).
		Where("email = ?", event.Sender.Email).Count(&lenEmail).Error
	if err != nil {
		glog.Errorf("failed to check user email: %v", err)
		return err
	}
	if lenEmail > 0 {
		// add label openeuler-cla/yes
		addlabel := &gitee.NoteEvent{}
		addlabel.PullRequest = event.PullRequest
		addlabel.Repository = event.Repository
		addlabel.Comment = &gitee.Note{}
		addlabel.Comment.Body = AddClaYes
		err = s.AddLabelInPulRequest(addlabel)
		if err != nil {
			return err
		}

		// remove label openeuler-cla/no
		removelabel := &gitee.NoteEvent{}
		removelabel.PullRequest = event.PullRequest
		removelabel.Repository = event.Repository
		removelabel.Comment = &gitee.Note{}
		removelabel.Comment.Body = RemoveClaNo
		err = s.RemoveLabelInPullRequest(removelabel)
		if err != nil {
			return err
		}
	} else {
		// add label openeuler-cla/no
		addlabel := &gitee.NoteEvent{}
		addlabel.PullRequest = event.PullRequest
		addlabel.Repository = event.Repository
		addlabel.Comment = &gitee.Note{}
		addlabel.Comment.Body = AddClaNo
		err = s.AddLabelInPulRequest(addlabel)
		if err != nil {
			return err
		}

		// remove label openeuler-cla/yes
		removelabel := &gitee.NoteEvent{}
		removelabel.PullRequest = event.PullRequest
		removelabel.Repository = event.Repository
		removelabel.Comment = &gitee.Note{}
		removelabel.Comment.Body = RemoveClaYes
		err = s.RemoveLabelInPullRequest(removelabel)
		if err != nil {
			return err
		}

		// add comment
		body := gitee.PullRequestCommentPostParam{}
		body.AccessToken = s.Config.GiteeToken
		body.Body = claNotFoundMessage
		owner := event.Repository.Owner.Login
		repo := event.Repository.Name
		number := event.PullRequest.Number
		_, _, err = s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(s.Context, owner, repo, number, body)
		if err != nil {
			glog.Errorf("unable to add comment in pull request: %v", err)
			return err
		}
	}
	return nil
}
