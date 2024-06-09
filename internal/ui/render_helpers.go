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

	t.desc = fmt.Sprintf("%s %s", RightPadTrim(lastUpdated, 60), timeSpent)
}

func (tl *taskLogEntry) updateDesc() {
	secsSpent := int(tl.endTS.Sub(tl.beginTS).Seconds())
	timeSpentStr := humanizeDuration(secsSpent)

	tl.desc = fmt.Sprintf("%s (spent %s)", RightPadTrim(humanize.Time(tl.beginTS), 60), timeSpentStr)
}
