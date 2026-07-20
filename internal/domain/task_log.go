package domain

import "time"

type TaskLogEntry struct {
	ID          int
	TaskID      int
	TaskSummary string
	BeginTS     time.Time
	EndTS       time.Time
	SecsSpent   int
	Comment     *string
}

type ActiveTaskDetails struct {
	TaskID            int
	TaskSummary       string
	CurrentLogBeginTS time.Time
	CurrentLogComment *string
}

type TaskReportEntry struct {
	TaskID      int
	TaskSummary string
	NumEntries  int
	SecsSpent   int
}
