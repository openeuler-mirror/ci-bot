package cibot

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	kind              = "/kind"
	RemoveKind        = "/remove-kind"
	AddClaYes         = "/openeuler-cla yes"
	AddClaNo          = "/openeuler-cla no"
	RemoveClaYes      = "/remove-openeuler-cla yes"
	RemoveClaNo       = "/remove-openeuler-cla no"
	LabelNameLgtm     = "lgtm"
	LabelNameApproved = "approved"
	LabelHiddenValue  = "<input type=hidden value=%s />"
	BotName           = "openeuler-ci-bot"
	tipBotMessage     = `Hey ***@%s***, Welcome to openEuler Community.
All of the projects in openEuler Community are maintained by ***@openeuler-ci-bot***.
That means the developpers can comment below every pull request or issue to trigger Bot Commands.
Please follow instructions at <https://gitee.com/openeuler/community/blob/master/en/command.md> to find the details.`
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
