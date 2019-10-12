package cibot

import (
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
)

// UrlEncode replcae special chars in url
func UrlEncode(str string) string {
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}
