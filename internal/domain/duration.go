package domain

import (
	"time"

	"github.com/dhth/hours/internal/types"
)

func SecondsTrackedToday(
	finishedTaskLogs []types.TaskLogEntry,
	activeTaskLogBeginTS *time.Time,
	now time.Time,
) int {
	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	var trackedSeconds int
	for _, entry := range finishedTaskLogs {
		trackedSeconds += overlappingSeconds(
			entry.BeginTS,
			entry.EndTS,
			startOfDay,
			now,
		)
	}

	if activeTaskLogBeginTS != nil && !activeTaskLogBeginTS.IsZero() {
		trackedSeconds += overlappingSeconds(
			*activeTaskLogBeginTS,
			now,
			startOfDay,
			now,
		)
	}

	return trackedSeconds
}

func overlappingSeconds(begin, end, rangeStart, rangeEnd time.Time) int {
	if !end.After(begin) || !rangeEnd.After(rangeStart) {
		return 0
	}

	overlapStart := begin
	if overlapStart.Before(rangeStart) {
		overlapStart = rangeStart
	}

	overlapEnd := end
	if overlapEnd.After(rangeEnd) {
		overlapEnd = rangeEnd
	}

	if !overlapEnd.After(overlapStart) {
		return 0
	}

	return int(overlapEnd.Sub(overlapStart).Seconds())
}
