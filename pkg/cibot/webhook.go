package cibot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
)

type Webhook struct {
	Address    string
	Port       int64
	ConfigFile string
}

func NewWebHook() *Webhook {
	return &Webhook{
		Address:    "0.0.0.0",
		Port:       10000,
		ConfigFile: "config.yaml",
	}
}

func (s *Webhook) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Address, "address", s.Address, "ip address to serve, 0.0.0.0 by default.")
	fs.Int64Var(&s.Port, "port", s.Port, "port to listen on, 10000 by default.")
	fs.StringVar(&s.ConfigFile, "config", s.ConfigFile, "config file.")
}

func (s *Webhook) Run() {
	// read file
	configContent, err := ioutil.ReadFile(s.ConfigFile)
	if err != nil {
		glog.Fatalf("could not read config file: %v", err)
	}

	// unmarshal config file
	var config Config
	err = json.Unmarshal(configContent, &config)
	if err != nil {
		glog.Fatalf("fail to unmarshal: %v", err)
	}

	// oauth
	oauthSecret := config.GiteeToken
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(oauthSecret)},
	)
	tc := oauth2.NewClient(ctx, ts)
	glog.Infof("oauth client: %v", tc)

	// return 200 for health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// setting handler
	webHookHandler := Server{
		Config:  config,
		Context: ctx,
	}
	http.HandleFunc("/hook", webHookHandler.ServeHTTP)

	//starting server
	address := s.Address + ":" + strconv.FormatInt(s.Port, 10)
	if err := http.ListenAndServe(address, nil); err != nil {
		glog.Error(err)
	}
}
