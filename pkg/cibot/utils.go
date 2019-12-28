package cibot

import (
	"fmt"
	"regexp"
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
)

const (
	kind              = "/kind"
	RemoveKind        = "/remove-kind"
	AddClaYes         = "/%s-cla yes"
	AddClaNo          = "/%s-cla no"
	RemoveClaYes      = "/remove-%s-cla yes"
	RemoveClaNo       = "/remove-%s-cla no"
	LabelNameLgtm     = "lgtm"
	LabelNameApproved = "approved"
	LabelHiddenValue  = "<input type=hidden value=%s />"
	tipBotMessage     = `Hey ***@%s***, Welcome to %s Community.
All of the projects in %s Community are maintained by ***@%s***.
That means the developpers can comment below every pull request or issue to trigger Bot Commands.
Please follow instructions at <%s> to find the details.`
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
)

// UrlEncode replcae special chars in url
func UrlEncode(str string) string {
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}

func canCommentPrIncludingSigDirectory(server *Server, owner string, repo string, prNumber int32, commentUser string) (int32, error) {
	files, _, err := server.GiteeClient.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberFiles(
		server.Context, owner, repo, prNumber, &gitee.GetV5ReposOwnerRepoPullsNumberFilesOpts{
			AccessToken: optional.NewString(server.Config.GiteeToken),
		})
	if err != nil {
		return -1, err
	}

	sigFilePathHeadPattern, err := regexp.Compile("^[a-zA-Z0-9_-]+/sig/[a-zA-Z0-9_-]+/")
	if err != nil {
		return -1, err
	}

	sigFilePathPattern, _ := regexp.Compile(sigFilePathHeadPattern.String() + ".+")

	targetSigPath := make(map[string]bool)

	for _, file := range files {
		// TODO test: use file.RawUrl instead?
		if !sigFilePathPattern.MatchString(file.Filename) {
			return -1, nil
		}

		targetSigPath[sigFilePathHeadPattern.FindString(file.Filename)] = true
	}

	for path, _ := range targetSigPath {
		content, _, err := server.GiteeClient.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
			server.Context, owner, repo, path+"OWNERS", &gitee.GetV5ReposOwnerRepoContentsPathOpts{
				AccessToken: optional.NewString(server.Config.GiteeToken),
			})
		// TODO the file is not exist. for example: thie pr create OWNERS.
		if err != nil {
			return -1, err
		}

		if !strings.Contains(content.Content, commentUser) {
			return 0, nil
		}
	}
	return 1, nil
}
