package ui

import "time"

type HideHelpMsg struct{}

type trackingToggledMsg struct {
	taskId    int
	finished  bool
	secsSpent int
	err       error
}

type taskRepUpdatedMsg struct {
	tsk *task
	err error
}

type manualTaskLogInserted struct {
	taskId int
	err    error
}

type activeTaskFetchedMsg struct {
	activeTaskId int
	beginTs      time.Time
	noneActive   bool
	err          error
}

type taskLogEntriesFetchedMsg struct {
	entries []taskLogEntry
	err     error
}

type taskCreatedMsg struct {
	err error
}

type taskUpdatedMsg struct {
	tsk     *task
	summary string
	err     error
}

type taskActiveStatusUpdated struct {
	tsk    *task
	active bool
	err    error
}

type taskLogEntryDeletedMsg struct {
	entry *taskLogEntry
	err   error
}

type tasksFetched struct {
	tasks  []task
	active bool
	err    error
}
