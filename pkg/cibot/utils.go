package cibot

import (
	"regexp"
	"strings"
)

const (
	kind         = "/kind"
	RemoveKind   = "/remove-kind"
	AddClaYes    = "/openeuler-cla yes"
	AddClaNo     = "/openeuler-cla no"
	RemoveClaYes = "/remove-openeuler-cla yes"
	RemoveClaNo  = "/remove-openeuler-cla no"
)

var (
	// RegAddLabel
	RegAddLabel = regexp.MustCompile(`(?mi)^/(kind|priority|sig|openeuler-cla)\s*(.*)$`)
	// RegRemoveLabel
	RegRemoveLabel = regexp.MustCompile(`(?mi)^/remove-(kind|priority|sig|openeuler-cla)\s*(.*)$`)
	// RegCheckCLA
	RegCheckCLA = regexp.MustCompile(`(?mi)^/check-cla\s*$`)
)

// UrlEncode replcae special chars in url
func UrlEncode(str string) string {
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}
