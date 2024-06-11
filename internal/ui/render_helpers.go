package ui

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func (t *task) updateTitle() {
	var trackingIndicator string
	if t.trackingActive {
		trackingIndicator = "⏲ "
	}

	t.title = trackingIndicator + t.summary
}

func (t *task) updateDesc() {
	var timeSpent string

	if t.secsSpent != 0 {
		timeSpent = "worked on for " + humanizeDuration(t.secsSpent)
	} else {
		timeSpent = "no time spent"
	}
	lastUpdated := fmt.Sprintf("last updated: %s", humanize.Time(t.updatedAt))

	t.desc = fmt.Sprintf("%s %s", RightPadTrim(lastUpdated, 60, true), timeSpent)
}

func (tl *taskLogEntry) updateDesc() {
	timeSpentStr := humanizeDuration(tl.secsSpent)

	timeStr := fmt.Sprintf("ended on %s (spent %s)",
		RightPadTrim(tl.endTS.Format(friendlyTimeFormat), 40, true),
		timeSpentStr)

	tl.desc = fmt.Sprintf("%s %s", RightPadTrim("["+tl.taskSummary+"]", 60, true), timeStr)
}
