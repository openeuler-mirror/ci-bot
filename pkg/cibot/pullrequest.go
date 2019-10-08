package cibot

import (
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

// HandlePullRequestEvent handles pull request event
func (s *Server) HandlePullRequestEvent(event *gitee.PullRequestEvent) {
	if event == nil {
		return
	}

	glog.Infof("pull request sendor: %v", event.Sender.Login)
}
