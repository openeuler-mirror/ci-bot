package cibot

import "regexp"

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
