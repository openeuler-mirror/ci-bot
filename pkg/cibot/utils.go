package cibot

import (
	"fmt"
	"regexp"
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"github.com/golang/glog"
)

const (
	kind                   = "/kind"
	RemoveKind             = "/remove-kind"
	AddClaYes              = "/%s-cla yes"
	AddClaNo               = "/%s-cla no"
	RemoveClaYes           = "/remove-%s-cla yes"
	RemoveClaNo            = "/remove-%s-cla no"
	LabelNameLgtm          = "lgtm"
	LabelLgtmWithCommenter = "lgtm-%s"
	LabelNameApproved      = "approved"
	LabelHiddenValue       = "<input type=hidden value=%s />"
	tipBotMessage          = `Hey ***%s***, Welcome to %s Community.
You can follow the instructions at <%s> to interact with the Bot.
%s`
	DisplayCommittors = `If you have any questions, you could contact SIG: [%s](%s), and maintainers: `
	SigPath           = `https://gitee.com/openeuler/community/tree/master/sig/%s`
	AutoAddPrjMsg = `Since you have added a item to the src-openeuler.yaml file, we will automatically generate a default package in project openEuler:Factory on OBS cluster for you.
If you need a more customized configuration, you can configure it according to the following instructions: `
)

var (
	// RegAddLabel
	RegAddLabel = regexp.MustCompile(`(?mi)^/(kind|priority|sig)\s*(.*)$`)
	// RegRemoveLabel
	RegRemoveLabel = regexp.MustCompile(`(?mi)^/remove-(kind|priority|sig)\s*(.*)$`)
	// RegCheckCLA
	RegCheckCLA = regexp.MustCompile(`(?mi)^/check-cla\s*$`)
	// RegAddLgtm
	RegAddLgtm = regexp.MustCompile(`(?mi)^/lgtm\s*$`)
	// RegRemoveLgtm
	RegRemoveLgtm = regexp.MustCompile(`(?mi)^/lgtm cancel\s*$`)
	// RegAddApprove
	RegAddApprove = regexp.MustCompile(`(?mi)^/approve\s*$`)
	// RegRemoveApprove
	RegRemoveApprove = regexp.MustCompile(`(?mi)^/approve cancel\s*$`)
	// RegClose
	RegClose = regexp.MustCompile(`(?mi)^/close\s*$`)
	// RegReOpen
	RegReOpen = regexp.MustCompile(`(?mi)^/reopen\s*$`)
	// RegBotAddLgtm
	RegBotAddLgtm = regexp.MustCompile(fmt.Sprintf(LabelHiddenValue, "(.*)"))
	// RegAssign
	RegAssign = regexp.MustCompile(`(?mi)^/assign(( @?[-\w]+?)*)\s*$`)
	// RegUnAssign
	RegUnAssign = regexp.MustCompile(`(?mi)^/unassign(( @?[-\w]+?)*)\s*$`)
	// RegCheckPr
	RegCheckPr = regexp.MustCompile(`(?mi)^/check-pr\s*$`)
)

// UrlEncode replcae special chars in url
func UrlEncode(str string) string {
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}

// canCommentPrIncludingSigDirectory
func canCommentPrIncludingSigDirectory(server *Server, owner string, repo string, prNumber int32, commentUser string) (int32, error) {
	// only check community
	if repo != "community" {
		return -1, nil
	}

	// get pr files
	files, _, err := server.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberFiles(
		server.Context, owner, repo, prNumber, &gitee.GetV5ReposOwnerRepoPullsNumberFilesOpts{
			AccessToken: optional.NewString(server.Config.GiteeToken),
		})
	if err != nil {
		glog.Infof("read pr files failed: %v", err)
		return -1, err
	}

	sigFilePathHeadPattern, err := regexp.Compile("^sig/[a-zA-Z0-9_-]+/")
	if err != nil {
		return -1, err
	}

	sigFilePathPattern, _ := regexp.Compile(sigFilePathHeadPattern.String() + ".+")
	targetSigPath := make(map[string]bool)

	for _, file := range files {
		// TODO test: use file.RawUrl instead?
		glog.Infof("get files log, RawUrl=%v, filename=%v\n", file.RawUrl, file.Filename)
		if !sigFilePathPattern.MatchString(file.Filename) {
			glog.Info("file name pattern log: not match")
			return -1, nil
		}
		targetSigPath[sigFilePathHeadPattern.FindString(file.Filename)] = true
	}

	glog.Infof("targetSigPath=%v", targetSigPath)

	for path, _ := range targetSigPath {
		content, _, err := server.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
			server.Context, owner, repo, path+"OWNERS", &gitee.GetV5ReposOwnerRepoContentsPathOpts{
				AccessToken: optional.NewString(server.Config.GiteeToken),
			})
		glog.Infof("read file=%v, err=%v", path, err)
		// TODO the file is not exist. for example: this pr create OWNERS.
		if err != nil {
			// TODO if read OWNERS failed, return nil to allow TC to approve this pr temporarily.
			return -1, nil
		}

		owners := DecodeOwners(content.Content)
		glog.Infof("owners=%v, commenturser=%v", owners, commentUser)
		if owners == nil {
			return 0, nil
		}

		bingo := false
		for _, owner := range owners {
			glog.Infof("owner=%v", owner)
			if owner == commentUser {
				bingo = true
				break
			}
		}
		if !bingo {
			return 0, nil
		}
	}
	return 1, nil
}

//truncateLabel on gitee the length of the label cannot exceed 20 characters.
//If it exceeds, the label will be truncated and replaced.
func truncateLabel(labels []string) []string {
	for i := range labels {
		rs := []rune(labels[i])
		if len(rs) > 20 {
			labels[i] = string(rs[:20])
		}
	}
	return labels
}
