package ui

import (
	"time"

	"github.com/dhth/hours/internal/types"
)

type hideHelpMsg struct{}

type trackingToggledMsg struct {
	taskID    int
	finished  bool
	secsSpent int
	err       error
}

type activeTLSwitchedMsg struct {
	lastActiveTaskID      int
	currentlyActiveTaskID int
	currentlyActiveTLID   int
	ts                    time.Time
	err                   error
}

type taskRepUpdatedMsg struct {
	tsk *types.Task
	err error
}

type manualTLInsertedMsg struct {
	taskID int
	err    error
}

type savedTLEditedMsg struct {
	tlID   int
	taskID int
	err    error
}

type activeTLUpdatedMsg struct {
	beginTS time.Time
	comment *string
	err     error
}

type activeTaskLogDeletedMsg struct {
	err error
}

type activeTaskFetchedMsg struct {
	activeTask types.ActiveTaskDetails
	noneActive bool
	err        error
}

type tLsFetchedMsg struct {
	entries       []types.TaskLogEntry
	tlIDToFocusOn *int
	err           error
}

type taskCreatedMsg struct {
	err error
}

type taskUpdatedMsg struct {
	tsk     *types.Task
	summary string
	err     error
}

type taskActiveStatusUpdatedMsg struct {
	tsk    *types.Task
	active bool
	err    error
}

type tLDeletedMsg struct {
	entry *types.TaskLogEntry
	err   error
}

type taskLogMovedMsg struct {
	tlID      int
	oldTaskID int
	newTaskID int
	err       error
}

type tasksFetchedMsg struct {
	tasks  []types.Task
	active bool
	err    error
}

type recordsDataFetchedMsg struct {
	dateRange types.DateRange
	report    string
	err       error
}
