package ui

import (
	"fmt"
	"time"

	"github.com/dhth/hours/internal/domain"
	"github.com/dhth/hours/internal/types"
)

type taskListItem struct {
	domain.Task
	trackingActive bool
	listTitle      string
	listDesc       string
}

func (t *taskListItem) updateListTitle() {
	var trackingIndicator string
	if t.trackingActive {
		trackingIndicator = "⏲ "
	}

	t.listTitle = trackingIndicator + t.Summary
}

// updateListDesc sets the task's list description. trackingSince is the begin
// time of the active task log when this task is being tracked, and nil otherwise.
func (t *taskListItem) updateListDesc(trackingSince *time.Time, now time.Time) {
	switch {
	case trackingSince == nil && t.SecsSpent == 0:
		t.listDesc = "no time spent"
	case trackingSince == nil:
		t.listDesc = "worked on for " + types.HumanizeDuration(t.SecsSpent)
	case t.SecsSpent == 0:
		t.listDesc = "currently tracking " + types.HumanizeDuration(int(now.Sub(*trackingSince).Seconds()))
	default:
		t.listDesc = fmt.Sprintf(
			"worked on for %s (currently tracking %s)",
			types.HumanizeDuration(t.SecsSpent),
			types.HumanizeDuration(int(now.Sub(*trackingSince).Seconds())),
		)
	}
}

func (t *taskListItem) Title() string {
	return t.listTitle
}

func (t *taskListItem) Description() string {
	return t.listDesc
}

func (t *taskListItem) FilterValue() string {
	return t.Summary
}
