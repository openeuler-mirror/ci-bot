package cibot

import (
	"fmt"
	"regexp"
	"strings"
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
