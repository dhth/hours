package ui

import (
	"fmt"

	"github.com/dhth/hours/internal/domain"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/utils"
	"github.com/dustin/go-humanize"
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

func (t *taskListItem) updateListDesc(timeProvider types.TimeProvider) {
	var timeSpent string
	if t.SecsSpent != 0 {
		timeSpent = "worked on for " + types.HumanizeDuration(t.SecsSpent)
	} else {
		timeSpent = "no time spent"
	}

	lastUpdated := fmt.Sprintf("last updated: %s", humanize.RelTime(t.UpdatedAt, timeProvider.Now(), "ago", "from now"))
	t.listDesc = fmt.Sprintf("%s %s", utils.RightPadTrim(lastUpdated, 60, true), timeSpent)
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
