package cibot

import (
	"context"
	"net/http"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

type Config struct {
	Owner         string `yaml:"owner"`
	Repo          string `yaml:"repository"`
	GiteeToken    string `yaml:"giteeToken"`
	WebhookSecret string `yaml:"webhookSecret"`
}

type Server struct {
	Config      Config
	Context     context.Context
	GiteeClient *gitee.APIClient
}

// ServeHTTP validates an incoming webhook and invoke its handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Info("received a webhook event")
	/* Validata the webhook secret
	payload, err := gitee.ValidatePayload(r, []byte(s.Config.WebhookSecret))
	if err != nil {
		glog.Errorf("Invalid payload: %v", err)
		return
	}*/

	/* Parse into Event
	event, err := gitee.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		glog.Errorf("Failed to parse webhook")
		return
	}*/
	glog.Infof("header: %v body: %v", r.Header, r.Body)

	var client http.Client
	client.Do(r)

	/* handle events
	switch event.(type) {
	case *gitee.IssueEvent:
		go s.handleIssueEvent(payload)
	}*/
}
