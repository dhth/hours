package ui

import (
	"time"
)

type task struct {
	id             int
	summary        string
	createdAt      time.Time
	updatedAt      time.Time
	trackingActive bool
	secsSpent      int
	active         bool
	title          string
	desc           string
}

type taskLogEntry struct {
	id          int
	taskId      int
	taskSummary string
	beginTs     time.Time
	endTs       time.Time
	secsSpent   int
	comment     string
	desc        string
}

type activeTaskDetails struct {
	taskId              int
	taskSummary         string
	lastLogEntryBeginTs time.Time
}

type taskReportEntry struct {
	taskId      int
	taskSummary string
	numEntries  int
	secsSpent   int
}

func (t task) Title() string {
	return Trim(t.title, 60)
}

func (t task) Description() string {
	return t.desc
}

func (t task) FilterValue() string {
	return t.summary
}

func (e taskLogEntry) Title() string {
	return Trim(e.comment, 60)
}

func (e taskLogEntry) Description() string {
	return e.desc
}

func (e taskLogEntry) FilterValue() string {
	return e.comment
}
