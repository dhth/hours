package ui

import (
	"time"

	"github.com/dhth/hours/internal/types"
)

type HideHelpMsg struct{}

type trackingToggledMsg struct {
	taskID    int
	finished  bool
	secsSpent int
	err       error
}

type taskRepUpdatedMsg struct {
	tsk *types.Task
	err error
}

type manualTaskLogInserted struct {
	taskID int
	err    error
}

type tlBeginTSUpdatedMsg struct {
	beginTS time.Time
	err     error
}

type activeTaskLogDeletedMsg struct {
	err error
}

type activeTaskFetchedMsg struct {
	activeTaskID int
	beginTs      time.Time
	noneActive   bool
	err          error
}

type taskLogEntriesFetchedMsg struct {
	entries []types.TaskLogEntry
	err     error
}

type taskCreatedMsg struct {
	err error
}

type taskUpdatedMsg struct {
	tsk     *types.Task
	summary string
	err     error
}

type taskActiveStatusUpdated struct {
	tsk    *types.Task
	active bool
	err    error
}

type taskLogEntryDeletedMsg struct {
	entry *types.TaskLogEntry
	err   error
}

type tasksFetched struct {
	tasks  []types.Task
	active bool
	err    error
}

type recordsDataFetchedMsg struct {
	start  time.Time
	end    time.Time
	report string
	err    error
}
