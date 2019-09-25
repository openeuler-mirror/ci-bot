package cibot

import (
	"context"
	"net/http"

	"gitee.com/openEuler/go-gitee/gitee"
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

}
