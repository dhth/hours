package ui

import (
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
	var timeSpent string
	if t.SecsSpent != 0 {
		timeSpent = "worked on for " + types.HumanizeDuration(t.SecsSpent)
	}

	if trackingSince == nil {
		if timeSpent == "" {
			timeSpent = "no time spent"
		}
		t.listDesc = timeSpent
		return
	}

	currentlyTracking := "currently tracking " + types.HumanizeDuration(int(now.Sub(*trackingSince).Seconds()))
	if timeSpent == "" {
		t.listDesc = currentlyTracking
		return
	}

	t.listDesc = timeSpent + " (" + currentlyTracking + ")"
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
