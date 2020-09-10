package cibot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

type Server struct {
	Config      config.Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

// ServeHTTP validates an incoming webhook and invoke its handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Info("received a webhook event")
	// validate the webhook secret
	payload, err := gitee.ValidatePayload(r, []byte(s.Config.WebhookSecret))
	if err != nil {
		glog.Errorf("invalid payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	// pstr := string(payload)
	// glog.Infof("payload: %v", pstr)

	// parse into Event
	messagetype := gitee.WebHookType(r)
	glog.Infof("message type: %v", messagetype)
	event, err := gitee.ParseWebHook(messagetype, payload)
	if err != nil {
		glog.Errorf("failed to parse webhook event: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	// response avoids gitee timeout 5s
	fmt.Fprint(w, "handle webhook event successfully")
	var client http.Client
	client.Do(r)

	// handle events
	switch event.(type) {
	case *gitee.NoteEvent:
		glog.Info("received a note event")
		go s.HandleNoteEvent(event.(*gitee.NoteEvent))
	case *gitee.PushEvent:
		glog.Info("received a push event")
		go s.HandlePushEvent(event.(*gitee.PushEvent))
	case *gitee.IssueEvent:
		glog.Info("received a issue event")
		go s.HandleIssueEvent(event.(*gitee.IssueEvent))
	case *gitee.PullRequestEvent:
		glog.Info("received a pull request event")
		actionDesc:= struct {
			ActionDesc string `json:"action_desc,omitempty"`
		}{}
		err := json.Unmarshal(payload, &actionDesc)
		if err != nil{
			glog.Info(err)
		}
		go s.HandlePullRequestEvent(actionDesc.ActionDesc,event.(*gitee.PullRequestEvent))
	case *gitee.TagPushEvent:
		glog.Info("received a tag push event")
	}
}
