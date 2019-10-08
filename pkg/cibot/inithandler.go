package cibot

import (
	"context"
	"io/ioutil"

	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

var (
	DefaultOrg = "openeuler"
)

type InitHandler struct {
	Config       config.Config
	Context      context.Context
	GiteeClient  *gitee.APIClient
	ProjectsFile string
}

type ProjectsFile struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Name        *string  `yaml:"name"`
	Owner       []string `yaml:"owner"`
	Type        *string  `yaml:"type"`
	Description *string  `yaml:"description"`
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

	glog.Infof("projects: %d", len(projects.Projects))
	for i := 0; i < len(projects.Projects); i++ {
		// log
		p := projects.Projects[i]
		glog.Infof("begin to run: %s", *p.Name)

		// build create repository param
		repobody := gitee.RepositoryPostParam{}
		repobody.AccessToken = handler.Config.GiteeToken
		repobody.Name = *p.Name
		repobody.Description = *p.Description
		repobody.HasIssues = true
		repobody.HasWiki = true
		if *p.Type == "private" {
			repobody.Private = true
		} else {
			glog.Infof("begin to public: %s", *p.Type)
			repobody.Private = false
		}

		// invoke create repository
		glog.Infof("begin to create repository: %s", *p.Name)
		_, _, err = handler.GiteeClient.RepositoriesApi.PostV5OrgsOrgRepos(handler.Context, DefaultOrg, repobody)
		if err != nil {
			glog.Errorf("fail to create repository: %v", err)
			continue
		}
		glog.Infof("end to create repository: %s", *p.Name)

		// build create project member param
		memberbody := gitee.ProjectMemberPutParam{}
		memberbody.AccessToken = handler.Config.GiteeToken
		memberbody.Permission = "admin"

		// invoke create project member
		glog.Infof("begin to create project member: %s", *p.Name)
		for j := 0; j < len(p.Owner); j++ {
			_, _, err = handler.GiteeClient.RepositoriesApi.PutV5ReposOwnerRepoCollaboratorsUsername(handler.Context, DefaultOrg, *p.Name, p.Owner[j], memberbody)
			if err != nil {
				glog.Errorf("fail to create project member: %v", err)
				continue
			}
		}
		glog.Infof("end to create project member: %s", *p.Name)

		// log
		glog.Infof("end to run: %s", *p.Name)
	}
}
