package ui

import (
	"fmt"
	"time"

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

func (tl *taskLogEntry) updateTitle() {
	tl.title = Trim(tl.comment, 60)
}

func (tl *taskLogEntry) updateDesc() {
	timeSpentStr := humanizeDuration(tl.secsSpent)

	var timeStr string
	var durationMsg string

	endTSRelative := getTSRelative(tl.endTs, time.Now())

	switch endTSRelative {
	case tsFromToday:
		durationMsg = fmt.Sprintf("%s  ...  %s", tl.beginTs.Format(timeOnlyFormat), tl.endTs.Format(timeOnlyFormat))
	case tsFromYesterday:
		durationMsg = "Yesterday"
	case tsFromThisWeek:
		durationMsg = tl.endTs.Format(dayFormat)
	default:
		durationMsg = humanize.Time(tl.endTs)
	}

	timeStr = fmt.Sprintf("%s (%s)",
		RightPadTrim(durationMsg, 40, true),
		timeSpentStr)

	tl.desc = fmt.Sprintf("%s %s", RightPadTrim("["+tl.taskSummary+"]", 60, true), timeStr)
}
