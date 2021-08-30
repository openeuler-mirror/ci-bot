package cibot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/antihax/optional"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

// GetLabelsMap for add or remove labels
func GetLabelsMap(comment string) map[string]string {
	// init labels map
	mapOfLabels := map[string]string{}
	// split with blank space
	substrings := strings.Split(strings.TrimSpace(comment), " ")
	// init label group
	labelGroup := ""
	// range over the substrings to get the map of labels
	for i, l := range substrings {
		if i == 0 {
			// first index is the operation to be performed, rest will be the labels
			// the label group. e.g kind, priority
			labelGroup = strings.Replace(strings.Replace(l, "/", "", 1), "remove-", "", 1)
		} else {
			// the whole label = label group + / + label. e.g kind/feature
			wholeLabel := labelGroup + "/" + l
			// use map to avoid the reduplicate label
			mapOfLabels[wholeLabel] = wholeLabel
		}
	}
	return mapOfLabels
}

//GetChangeLabels return the exact list of delete and update labels
func GetChangeLabels(confDelLabels []string, prLabels []gitee.Label) (delLabels, updateLabels []string) {
	for _, pl := range prLabels {
		isUpdate := true
		for _, cdl := range confDelLabels {
			//handle lgtm label
			if cdl == LabelNameLgtm {
				if strings.HasPrefix(pl.Name, fmt.Sprintf(LabelLgtmWithCommenter, "")) || pl.Name == LabelNameLgtm {
					delLabels = append(delLabels, pl.Name)
					isUpdate = false
				}
				continue
			}
			if cdl == "openeuler-cla" {
				if strings.HasPrefix(pl.Name, cdl) {
					delLabels = append(delLabels, pl.Name)
					isUpdate = false
				}
				continue
			}
			if cdl == "sig" {
				if strings.HasPrefix(pl.Name, cdl) {
					delLabels = append(delLabels, pl.Name)
					isUpdate = false
				}
				continue
			}
			if cdl == "kind" {
				if strings.HasPrefix(pl.Name, cdl) {
					delLabels = append(delLabels, pl.Name)
					isUpdate = false
				}
				continue
			}
			if cdl == pl.Name {
				delLabels = append(delLabels, pl.Name)
				isUpdate = false
				continue
			}

		}
		if isUpdate {
			updateLabels = append(updateLabels, pl.Name)
		}
	}
	return
}

// GetListOfAddLabels return the exact list of add labels
func GetListOfAddLabels(mapOfAddLabels map[string]string, listofRepoLabels []gitee.Label, listofItemLabels []gitee.Label) []string {
	// init
	listOfAddLabels := make([]string, 0)
	// range over the map to filter the list of labels
	for l := range mapOfAddLabels {
		// check if the label is existing in current gitee repository
		existingInRepo := false
		for _, repoLabel := range listofRepoLabels {
			if l == repoLabel.Name {
				existingInRepo = true
				break
			}
		}
		// the label is not existing in current gitee repository so it can not add this label
		if !existingInRepo {
			glog.Infof("label %s is not existing in repository", l)
			continue
		}

		// check if the label is existing in current item
		existingInItem := false
		for _, itemLabel := range listofItemLabels {
			if l == itemLabel.Name {
				existingInItem = true
				break
			}
		}
		// the label is already existing in current item so it is no need to add this label
		if existingInItem {
			glog.Infof("label %s is already existing in current item", l)
			continue
		}

		// append
		listOfAddLabels = append(listOfAddLabels, l)
	}
	return listOfAddLabels
}

// GetListOfRemoveLabels return the exact list of remove labels
func GetListOfRemoveLabels(mapOfRemoveLabels map[string]string, listofItemLabels []gitee.Label) []string {
	// init
	listOfRemoveLabels := make([]string, 0)
	// range over the map to filter the list of labels
	for l := range mapOfRemoveLabels {
		// check if the label is existing in current item
		existingInItem := false
		for _, itemLabel := range listofItemLabels {
			if l == itemLabel.Name {
				existingInItem = true
				break
			}
		}
		// the label is not existing in current item so it is no need to remove this label
		if !existingInItem {
			glog.Infof("label %s is not existing in current item", l)
			continue
		}

		// append
		listOfRemoveLabels = append(listOfRemoveLabels, l)
	}
	return listOfRemoveLabels
}

// AddLabel adds label
func (s *Server) AddLabel(event *gitee.NoteEvent) error {
	if *event.NoteableType == "PullRequest" {
		// PullRequest
		return s.AddLabelInPulRequest(event)
	} else if *event.NoteableType == "Issue" {
		// Issue
		return s.AddLabelInIssue(event)
	} else {
		return nil
	}
}

// AddLabelInPulRequest adds label in pull request
func (s *Server) AddLabelInPulRequest(event *gitee.NoteEvent) error {
	// get basic informations
	comment := event.Comment.Body
	owner := event.Repository.Namespace
	repo := event.Repository.Path
	var number int32
	if event.PullRequest != nil {
		number = event.PullRequest.Number
	}
	glog.Infof("add label started. comment: %s owner: %s repo: %s number: %d",
		comment, owner, repo, number)

	// /kind label1
	// /kind label2
	getLabels := strings.Split(comment, "\r\n")

	for _, labelToAdd := range getLabels {
		// map of add labels
		mapOfAddLabels := GetLabelsMap(labelToAdd)
		glog.Infof("map of add labels: %v", mapOfAddLabels)

		// list labels in current gitee repository
		lvosRepo := &gitee.GetV5ReposOwnerRepoLabelsOpts{}
		lvosRepo.AccessToken = optional.NewString(s.Config.GiteeToken)
		listofRepoLabels, _, err := s.GiteeClient.LabelsApi.GetV5ReposOwnerRepoLabels(s.Context, owner, repo, lvosRepo)
		if err != nil {
			glog.Errorf("unable to list repository labels. err: %v", err)
			return err
		}
		glog.Infof("list of repository labels: %v", listofRepoLabels)

		// list labels in current item
		lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{}
		lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
		pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, lvos)
		if err != nil {
			glog.Errorf("unable to get pull request. err: %v", err)
			return err
		}
		listofItemLabels := pr.Labels
		glog.Infof("list of item labels: %v", listofItemLabels)

		// list of add labels
		listOfAddLabels := GetListOfAddLabels(mapOfAddLabels, listofRepoLabels, listofItemLabels)
		glog.Infof("list of add labels: %v", listOfAddLabels)

		// invoke gitee api to add labels
		if len(listOfAddLabels) > 0 {
			// build label string
			var strLabel string
			for _, currentlabel := range listofItemLabels {
				strLabel += currentlabel.Name + ","
			}
			for _, addedlabel := range listOfAddLabels {
				strLabel += addedlabel + ","
			}
			strLabel = strings.TrimRight(strLabel, ",")
			body := gitee.PullRequestUpdateParam{}
			body.AccessToken = s.Config.GiteeToken
			body.Labels = strLabel
			glog.Infof("invoke api to add labels: %v", strLabel)

			// patch labels
			_, response, err := s.GiteeClient.PullRequestsApi.PatchV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, body)
			if err != nil {
				if response.StatusCode == 400 {
					glog.Infof("add labels successfully with status code %d: %v", response.StatusCode, listOfAddLabels)
				} else {
					glog.Errorf("unable to add labels: %v err: %v", listOfAddLabels, err)
					return err
				}
			} else {
				glog.Infof("add labels successfully: %v", listOfAddLabels)
			}
		} else {
			glog.Infof("no label to add for this event")
		}
	}

	return nil
}

func (s *Server) labelDiffer(new []string, existing []gitee.Label) ([]string, []string) {
	mb := make(map[string]struct{}, len(existing))
	for _, x := range existing {
		mb[x.Name] = struct{}{}
	}
	var diff []string
	var same []string
	for _, x := range new {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		} else {
			same = append(same, x)
		}
	}
	return diff, same
}

// AddLabelInIssue adds label in issue
func (s *Server) AddLabelInIssue(event *gitee.NoteEvent) error {
	// get basic informations
	comment := event.Comment.Body
	owner := event.Repository.Namespace
	repo := event.Repository.Path
	var number string
	if event.Issue != nil {
		number = event.Issue.Number
	}
	glog.Infof("add label started. comment: %s owner: %s repo: %s number: %s",
		comment, owner, repo, number)

	// /kind label1
	// /kind label2
	getLabels := strings.Split(comment, "\r\n")

	for _, labelToAdd := range getLabels {
		// map of add labels
		mapOfAddLabels := GetLabelsMap(labelToAdd)
		glog.Infof("map of add labels: %v", mapOfAddLabels)

		// list labels in current gitee repository
		lvosRepo := &gitee.GetV5ReposOwnerRepoLabelsOpts{}
		lvosRepo.AccessToken = optional.NewString(s.Config.GiteeToken)
		listofRepoLabels, _, err := s.GiteeClient.LabelsApi.GetV5ReposOwnerRepoLabels(s.Context, owner, repo, lvosRepo)
		if err != nil {
			glog.Errorf("unable to list repository labels. err: %v", err)
			return err
		}
		glog.Infof("list of repository labels: %v", listofRepoLabels)

		// list labels in current item
		lvos := &gitee.GetV5ReposOwnerRepoIssuesNumberLabelsOpts{}
		lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
		listofItemLabels, _, err := s.GiteeClient.LabelsApi.GetV5ReposOwnerRepoIssuesNumberLabels(s.Context, owner, repo, number, lvos)
		if err != nil {
			glog.Errorf("unable to get labels in issue. err: %v", err)
			return err
		}
		glog.Infof("list of item labels: %v", listofItemLabels)

		// list of add labels
		listOfAddLabels := GetListOfAddLabels(mapOfAddLabels, listofRepoLabels, listofItemLabels)
		glog.Infof("list of add labels: %v", listOfAddLabels)

		// invoke gitee api to add labels
		if len(listOfAddLabels) > 0 {
			// build label string
			var strLabel string
			for _, currentlabel := range listofItemLabels {
				strLabel += currentlabel.Name + ","
			}
			for _, addedlabel := range listOfAddLabels {
				strLabel += addedlabel + ","
			}
			strLabel = strings.TrimRight(strLabel, ",")
			body := gitee.IssueUpdateParam{}
			body.Repo = repo
			body.AccessToken = s.Config.GiteeToken
			body.Labels = strLabel
			glog.Infof("invoke api to add labels: %v", strLabel)

			// patch labels
			_, _, err := s.GiteeClient.IssuesApi.PatchV5ReposOwnerIssuesNumber(s.Context, owner, number, body)
			if err != nil {
				glog.Errorf("unable to add labels: %v err: %v", listOfAddLabels, err)
				return err
			} else {
				glog.Infof("add labels successfully: %v", listOfAddLabels)
			}
		} else {
			glog.Infof("no label to add for this event")
		}
	}

	return nil
}

// RemoveLabel removes label
func (s *Server) RemoveLabel(event *gitee.NoteEvent) error {
	if *event.NoteableType == "PullRequest" {
		// PullRequest
		return s.RemoveLabelInPullRequest(event)
	} else if *event.NoteableType == "Issue" {
		// Issue
		return s.RemoveLabelInIssue(event)
	} else {
		return nil
	}
}

// RemoveLabelInPullRequest removes label in pull request
func (s *Server) RemoveLabelInPullRequest(event *gitee.NoteEvent) error {
	// get basic informations
	comment := event.Comment.Body
	owner := event.Repository.Namespace
	repo := event.Repository.Path
	var number int32
	if event.PullRequest != nil {
		number = event.PullRequest.Number
	}
	glog.Infof("remove label started. comment: %s owner: %s repo: %s number: %d",
		comment, owner, repo, number)

	// /remove-kind label1
	// /remove-kind label2
	getLables := strings.Split(comment, "\r\n")

	for _, labelToRemove := range getLables {
		// map of add labels
		mapOfRemoveLabels := GetLabelsMap(labelToRemove)
		glog.Infof("map of remove labels: %v", mapOfRemoveLabels)

		// list labels in current item
		lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{}
		lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
		pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, lvos)
		if err != nil {
			glog.Errorf("unable to get pull request. err: %v", err)
			return err
		}
		listofItemLabels := pr.Labels
		glog.Infof("list of item labels: %v", listofItemLabels)

		// list of remove labels
		listOfRemoveLabels := GetListOfRemoveLabels(mapOfRemoveLabels, listofItemLabels)
		glog.Infof("list of remove labels: %v", listOfRemoveLabels)

		// invoke gitee api to remove labels
		if len(listOfRemoveLabels) > 0 {
			// build label string
			var strLabel string
			for _, currentlabel := range listofItemLabels {
				strLabel += currentlabel.Name + ","
			}
			for _, removedlabel := range listOfRemoveLabels {
				strLabel = strings.Replace(strLabel, removedlabel+",", "", 1)

			}
			strLabel = strings.TrimRight(strLabel, ",")
			// avoid to unable to remove labels when no label is exsit
			if strLabel == "" {
				strLabel = ","
			}
			body := gitee.PullRequestUpdateParam{}
			body.AccessToken = s.Config.GiteeToken
			body.Labels = strLabel
			glog.Infof("invoke api to remove labels: %v", strLabel)

			// patch labels
			_, response, err := s.GiteeClient.PullRequestsApi.PatchV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, body)
			if err != nil {
				if response.StatusCode == 400 {
					glog.Infof("remove labels successfully with status code %d: %v", response.StatusCode, listOfRemoveLabels)
				} else {
					glog.Errorf("unable to remove labels: %v err: %v", listOfRemoveLabels, err)
					return err
				}
			} else {
				glog.Infof("remove labels successfully: %v", listOfRemoveLabels)
			}
		} else {
			glog.Infof("no label to remove for this event")
		}
	}

	return nil
}

// RemoveLabelInIssue removes label in issue
func (s *Server) RemoveLabelInIssue(event *gitee.NoteEvent) error {
	// get basic informations
	comment := event.Comment.Body
	owner := event.Repository.Namespace
	repo := event.Repository.Path
	var number string
	if event.Issue != nil {
		number = event.Issue.Number
	}
	glog.Infof("remove label started. comment: %s owner: %s repo: %s number: %s",
		comment, owner, repo, number)

	// /remove-kind label1
	// /remove-kind label2
	getLables := strings.Split(comment, "\r\n")

	for _, labelToRemove := range getLables {
		// map of add labels
		mapOfRemoveLabels := GetLabelsMap(labelToRemove)
		glog.Infof("map of remove labels: %v", mapOfRemoveLabels)

		// list labels in current item
		lvos := &gitee.GetV5ReposOwnerRepoIssuesNumberLabelsOpts{}
		lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
		listofItemLabels, _, err := s.GiteeClient.LabelsApi.GetV5ReposOwnerRepoIssuesNumberLabels(s.Context, owner, repo, number, lvos)
		if err != nil {
			glog.Errorf("unable to get labels in issue. err: %v", err)
			return err
		}
		glog.Infof("list of item labels: %v", listofItemLabels)

		// list of remove labels
		listOfRemoveLabels := GetListOfRemoveLabels(mapOfRemoveLabels, listofItemLabels)
		glog.Infof("list of remove labels: %v", listOfRemoveLabels)

		// invoke gitee api to remove labels
		if len(listOfRemoveLabels) > 0 {
			glog.Infof("invoke api to remove labels: %v", listOfRemoveLabels)
			// remove labels
			for _, removedlabel := range listOfRemoveLabels {
				localVarOptionals := &gitee.DeleteV5ReposOwnerRepoIssuesNumberLabelsNameOpts{}
				localVarOptionals.AccessToken = optional.NewString(s.Config.GiteeToken)
				_, err := s.GiteeClient.LabelsApi.DeleteV5ReposOwnerRepoIssuesNumberLabelsName(
					s.Context, owner, repo, number, UrlEncode(removedlabel), localVarOptionals)
				if err != nil {
					glog.Errorf("unable to remove label: %s err: %v", removedlabel, err)
					return err
				} else {
					glog.Infof("remove label successfully: %s", removedlabel)
				}
			}
		} else {
			glog.Infof("no label to remove for this event")
		}
	}

	return nil
}

func (s *Server) generateLabelLuckyColor() string {
	cRand := rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%02x%02x%02x", cRand.Intn(255), cRand.Intn(255), cRand.Intn(255))
}

func (s *Server) patchPRLabels(labels []string, owner, repo string, number int32) error {
	// list labels in current pr
	lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{
		AccessToken: optional.NewString(s.Config.GiteeToken),
	}
	pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, lvos)
	if err != nil {
		glog.Errorf("unable to get pull request. err: %v", err)
		return err
	}
	glog.Infof("list of pr labels: %v", pr.Labels)

	// list of add labels
	newLabels, _ := s.labelDiffer(labels, pr.Labels)
	if len(newLabels) == 0 {
		glog.Info("all labels existed in pr, skip updating.")
		return nil
	}
	glog.Infof("list of add labels in pr: %v", newLabels)

	// invoke gitee api to add labels
	if len(newLabels) > 0 {
		// build label string
		body := gitee.PullRequestLabelPostParam{}
		body.AccessToken = s.Config.GiteeToken
		body.Body = newLabels
		glog.Infof("invoke api to add labels: %v", newLabels)

		// patch labels
		_, response, err := s.GiteeClient.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberLabels(s.Context, owner, repo, number, body)
		if err != nil {
			if response.StatusCode == 400 {
				glog.Infof("add labels successfully with status code %d: %v", response.StatusCode, newLabels)
			} else {
				glog.Errorf("unable to add labels: %v err: %v", newLabels, err)
				return err
			}
		} else {
			glog.Infof("add labels successfully: %v", newLabels)
		}
	} else {
		glog.Infof("no label to add for this event")
	}
	return nil
}

func (s *Server) patchRepoLabels(labels []string, owner, repo string, createNew bool) ([]string, error) {
	// list labels in current gitee repository
	lvosRepo := &gitee.GetV5ReposOwnerRepoLabelsOpts{
		AccessToken: optional.NewString(s.Config.GiteeToken),
	}
	listofRepoLabels, _, err := s.GiteeClient.LabelsApi.GetV5ReposOwnerRepoLabels(s.Context, owner, repo, lvosRepo)
	if err != nil {
		glog.Errorf("unable to list repository labels. err: %v", err)
		return []string{}, err
	}
	glog.Infof("list of repository labels: %v", listofRepoLabels)
	newLabels, existlabel := s.labelDiffer(labels, listofRepoLabels)
	if len(newLabels) == 0 {
		glog.Info("all labels existed in repo, skip creating.")
		return existlabel, nil
	}
	if !createNew {
		glog.Info("'createNew' flag is false, skip creating.")
		return existlabel, nil
	}
	createLabelParam := gitee.LabelPostParam{
		AccessToken: s.Config.GiteeToken,
	}
	for _, label := range newLabels {
		createLabelParam.Name = label
		createLabelParam.Color = s.generateLabelLuckyColor()
		result, httpResponse, err := s.GiteeClient.LabelsApi.PostV5ReposOwnerRepoLabels(
			s.Context, owner, repo, createLabelParam)
		if err != nil {
			glog.Errorf("unable to create repository label. %v, err: %v and %v", result, err.Error(), *httpResponse)
			return []string{}, err
		}
	}
	return labels, nil
}

// AddSpecifyLabelInPulRequest adds specify labels in pull request
func (s *Server) AddSpecifyLabelsInPulRequest(event *gitee.NoteEvent, newLabels []string, createNew bool) error {
	newLabels = truncateLabel(newLabels)
	// get basic information
	comment := event.Comment.Body
	owner := event.Repository.Namespace
	repo := event.Repository.Path
	var number int32
	if event.PullRequest != nil {
		number = event.PullRequest.Number
	}
	glog.Infof("add specify label started. comment: %s owner: %s repo: %s number: %d",
		comment, owner, repo, number)

	// patch repo labels
	legalLabels, err := s.patchRepoLabels(newLabels, owner, repo, createNew)
	if err != nil {
		return err
	}

	// patch pr labels
	if err := s.patchPRLabels(legalLabels, owner, repo, number); err != nil {
		return err
	}
	return nil
}

// RemoveSpecifyLabelsInPulRequest removes specify labels in pull request
func (s *Server) RemoveSpecifyLabelsInPulRequest(event *gitee.NoteEvent, mapOfRemoveLabels map[string]string) error {
	// get basic information
	comment := event.Comment.Body
	owner := event.Repository.Namespace
	repo := event.Repository.Path
	var number int32
	if event.PullRequest != nil {
		number = event.PullRequest.Number
	}
	glog.Infof("remove label started. comment: %s owner: %s repo: %s number: %d",
		comment, owner, repo, number)

	// list labels in current item
	lvos := &gitee.GetV5ReposOwnerRepoPullsNumberOpts{}
	lvos.AccessToken = optional.NewString(s.Config.GiteeToken)
	pr, _, err := s.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, lvos)
	if err != nil {
		glog.Errorf("unable to get pull request. err: %v", err)
		return err
	}
	listofItemLabels := pr.Labels
	glog.Infof("list of item labels: %v", listofItemLabels)

	// list of remove labels
	listOfRemoveLabels := GetListOfRemoveLabels(mapOfRemoveLabels, listofItemLabels)
	glog.Infof("list of remove labels: %v", listOfRemoveLabels)

	// invoke gitee api to remove labels
	if len(listOfRemoveLabels) > 0 {
		// build label string
		var strLabel string
		for _, currentlabel := range listofItemLabels {
			strLabel += currentlabel.Name + ","
		}
		for _, removedlabel := range listOfRemoveLabels {
			strLabel = strings.Replace(strLabel, removedlabel+",", "", 1)

		}
		strLabel = strings.TrimRight(strLabel, ",")
		// avoid to unable to remove labels when no label is exsit
		if strLabel == "" {
			strLabel = ","
		}
		body := gitee.PullRequestUpdateParam{}
		body.AccessToken = s.Config.GiteeToken
		body.Labels = strLabel
		glog.Infof("invoke api to remove labels: %v", strLabel)

		// patch labels
		_, response, err := s.GiteeClient.PullRequestsApi.PatchV5ReposOwnerRepoPullsNumber(s.Context, owner, repo, number, body)
		if err != nil {
			if response.StatusCode == 400 {
				glog.Infof("remove labels successfully with status code %d: %v", response.StatusCode, listOfRemoveLabels)
			} else {
				glog.Errorf("unable to remove labels: %v err: %v", listOfRemoveLabels, err)
				return err
			}
		} else {
			glog.Infof("remove labels successfully: %v", listOfRemoveLabels)
		}
	} else {
		glog.Infof("no label to remove for this event")
	}

	return nil
}
