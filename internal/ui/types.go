package ui

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

type task struct {
	id             int
	summary        string
	createdAt      time.Time
	updatedAt      time.Time
	trackingActive bool
	title          string
	desc           string
	secsSpent      int
}

type taskLogEntry struct {
	id          int
	taskId      int
	taskSummary string
	beginTS     time.Time
	endTS       time.Time
	comment     string
}

func (t task) Title() string {
	return t.title
}

func (t task) Description() string {
	return t.desc
}

func (t task) FilterValue() string {
	return t.summary
}

func (e taskLogEntry) Title() string {
	return e.comment
}

func (e taskLogEntry) Description() string {
	secsSpent := int(e.endTS.Sub(e.beginTS).Seconds())
	timeSpentStr := humanizeDuration(secsSpent)

	timeStr := fmt.Sprintf("began %s (spent %s)", RightPadTrim(humanize.Time(e.beginTS), 20), timeSpentStr)

	return fmt.Sprintf("%s %s", RightPadTrim("["+e.taskSummary+"]", 60), timeStr)
}

func (e taskLogEntry) FilterValue() string {
	return e.comment
}
