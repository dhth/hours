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
		timeSpent = "time spent: " + humanizeDuration(t.secsSpent)
	} else {
		timeSpent = "no time spent"
	}
	lastUpdated := fmt.Sprintf("last updated: %s", humanize.Time(t.updatedAt))

	t.desc = fmt.Sprintf("%s %s", RightPadTrim(lastUpdated, 60), timeSpent)
}
