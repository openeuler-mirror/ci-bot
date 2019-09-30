package cibot

import (
	"context"
	"io/ioutil"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

type InitHandler struct {
	Config       Config
	Context      context.Context
	GiteeClient  *gitee.APIClient
	ProjectsFile string
}

type ProjectsFile struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Name        string `yaml:"name"`
	Owner       string `yaml:"owner"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

// Serve
func (handler *InitHandler) Serve() {
	// read file
	projectsContent, err := ioutil.ReadFile(handler.ProjectsFile)
	if err != nil {
		glog.Errorf("could not read projects file: %v", err)
	}

	// unmarshal projects file
	var projects ProjectsFile
	err = yaml.Unmarshal(projectsContent, &projects)
	if err != nil {
		glog.Errorf("fail to unmarshal: %v", err)
	}

	glog.Infof("projects: %v", projects)
}
