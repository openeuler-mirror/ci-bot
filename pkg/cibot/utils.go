package cibot

import (
	"regexp"
	"strings"
)

const (
	kind       = "/kind"
	RemoveKind = "/remove-kind"
)

var (
	// RegAddLabel
	RegAddLabel = regexp.MustCompile(`(?mi)^/(kind|priority)\s*(.*)$`)
	// RegRemoveLabel
	RegRemoveLabel = regexp.MustCompile(`(?mi)^/remove-(kind|priority)\s*(.*)$`)
)

// UrlEncode replcae special chars in url
func UrlEncode(str string) string {
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}
