package cibot

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"

	cfg "gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Webhook struct {
	Address    string
	Port       int64
	ConfigFile string
}

func NewWebHook() *Webhook {
	return &Webhook{
		Address:    "0.0.0.0",
		Port:       8888,
		ConfigFile: "config.yaml",
	}
}

func (s *Webhook) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Address, "address", s.Address, "ip address to serve, 0.0.0.0 by default.")
	fs.Int64Var(&s.Port, "port", s.Port, "port to listen on, 8888 by default.")
	fs.StringVar(&s.ConfigFile, "configfile", s.ConfigFile, "config file.")

	// See https://github.com/spf13/pflag#supporting-go-flags-when-using-pflag
	fs.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
}

func (s *Webhook) Run() {
	// Flush flushes all pending log I/O.
	defer glog.Flush()

	// read file
	configContent, err := ioutil.ReadFile(s.ConfigFile)
	if err != nil {
		glog.Fatalf("could not read config file: %v", err)
	}

	// unmarshal config file
	var config cfg.Config
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		glog.Fatalf("fail to unmarshal: %v", err)
	}

	//parse environment variables by tag
	err = cfg.ParseEnvConf(&config, "")
	if err != nil {
    	glog.Info("fail to ParseEnvConf: %v",err)
	}

	// oauth
	oauthSecret := config.GiteeToken
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(oauthSecret)},
	)

	// configuration
	giteeConf := gitee.NewConfiguration()
	giteeConf.HTTPClient = oauth2.NewClient(ctx, ts)

	// git client
	giteeClient := gitee.NewAPIClient(giteeConf)

	err = database.New(config)
	if err != nil {
		glog.Errorf("init back database error: %v", err)
	}

	/* setting init handler
	initHandler := InitHandler{
		Config:      config,
		Context:     ctx,
		GiteeClient: giteeClient,
	}
	go initHandler.Serve()*/

	// setting repo handler
	repoHandler := RepoHandler{
		Config:      config,
		Context:     ctx,
		GiteeClient: giteeClient,
	}
	go repoHandler.Serve()

	// setting sig handler
	sigHandler := SigHandler{
		Config:      config,
		Context:     ctx,
		GiteeClient: giteeClient,
	}
	go sigHandler.Serve()

	// setting owner handler
	ownerHandler := OwnerHandler{
		Config:      config,
		Context:     ctx,
		GiteeClient: giteeClient,
	}
	go ownerHandler.Serve()

	// return 200 for health check
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	// setting webhook handler
	webHookHandler := Server{
		Config:      config,
		Context:     ctx,
		GiteeClient: giteeClient,
	}
	http.HandleFunc("/webhook", webHookHandler.ServeHTTP)

	// setting cla handler
	claHandler := CLAHandler{
		Context: ctx,
	}
	http.HandleFunc("/cla", claHandler.ServeHTTP)

	//starting server
	address := s.Address + ":" + strconv.FormatInt(s.Port, 10)
	if err := http.ListenAndServe(address, nil); err != nil {
		glog.Error(err)
	}
}
