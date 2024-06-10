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
	beginTS     time.Time
	endTS       time.Time
	comment     string
	title       string
	desc        string
}

type activeTaskDetails struct {
	taskId              int
	taskSummary         string
	lastLogEntryBeginTS time.Time
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
	return e.desc
}

func (e taskLogEntry) FilterValue() string {
	return e.comment
}
