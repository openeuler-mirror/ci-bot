package cibot

import (
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
)

// HandleNoteEvent handles note event
func (s *Server) HandleNoteEvent(event *gitee.NoteEvent) {
	if event == nil {
		return
	}
	// just handle create comment event
	if *event.Action != "comment" {
		return
	}

	// add label
	if RegAddLabel.MatchString(event.Comment.Body) {
		err := s.AddLabel(event)
		if err != nil {
			glog.Errorf("failed to add label: %v", err)
		}
	}

	// remove label
	if RegRemoveLabel.MatchString(event.Comment.Body) {
		err := s.RemoveLabel(event)
		if err != nil {
			glog.Errorf("failed to remove label: %v", err)
		}
	}

	// check cla by note event
	if RegCheckCLA.MatchString(event.Comment.Body) {
		err := s.CheckCLAByNoteEvent(event)
		if err != nil {
			glog.Errorf("failed to check cla by note event: %v", err)
		}
	}

	// add lgtm
	if RegAddLgtm.MatchString(event.Comment.Body) {
		err := s.AddLgtm(event)
		if err != nil {
			glog.Errorf("failed to add lgtm: %v", err)
		}
	}

	// remove lgtm
	if RegRemoveLgtm.MatchString(event.Comment.Body) {
		err := s.RemoveLgtm(event)
		if err != nil {
			glog.Errorf("failed to remove lgtm: %v", err)
		}
	}

	// add approve
	if RegAddApprove.MatchString(event.Comment.Body) {
		err := s.AddApprove(event)
		if err != nil {
			glog.Errorf("failed to add approved: %v", err)
		}
	}

	// remove approve
	if RegRemoveApprove.MatchString(event.Comment.Body) {
		err := s.RemoveApprove(event)
		if err != nil {
			glog.Errorf("failed to remove approved: %v", err)
		}
	}

	// close
	if RegClose.MatchString(event.Comment.Body) {
		err := s.Close(event)
		if err != nil {
			glog.Errorf("failed to close: %v", err)
		}
	}

	// reopen
	if RegReOpen.MatchString(event.Comment.Body) {
		err := s.ReOpen(event)
		if err != nil {
			glog.Errorf("failed to reopen: %v", err)
		}
	}

	// assign
	if RegAssign.MatchString(event.Comment.Body) {
		err := s.Assign(event)
		if err != nil {
			glog.Errorf("failed to assign: %v", err)
		}
	}

	// unassign
	if RegUnAssign.MatchString(event.Comment.Body) {
		err := s.UnAssign(event)
		if err != nil {
			glog.Errorf("failed to unassign: %v", err)
		}
	}
}
