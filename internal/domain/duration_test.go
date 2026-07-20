package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverlappingSeconds(t *testing.T) {
	// GIVEN
	rangeStart := timestamp(t, "2026-07-16T09:00:00Z")
	rangeEnd := timestamp(t, "2026-07-16T17:00:00Z")
	testCases := []struct {
		name            string
		begin           string
		end             string
		expectedSeconds int
	}{
		{
			name:            "entirely within range",
			begin:           "2026-07-16T10:00:00Z",
			end:             "2026-07-16T12:30:00Z",
			expectedSeconds: (2*60 + 30) * 60,
		},
		{
			name:            "overlaps range start",
			begin:           "2026-07-16T08:00:00Z",
			end:             "2026-07-16T10:00:00Z",
			expectedSeconds: 60 * 60,
		},
		{
			name:            "overlaps range end",
			begin:           "2026-07-16T16:00:00Z",
			end:             "2026-07-16T18:00:00Z",
			expectedSeconds: 60 * 60,
		},
		{
			name:            "spans entire range",
			begin:           "2026-07-16T08:00:00Z",
			end:             "2026-07-16T18:00:00Z",
			expectedSeconds: 8 * 60 * 60,
		},
		{
			name:            "entirely before range",
			begin:           "2026-07-16T07:00:00Z",
			end:             "2026-07-16T08:00:00Z",
			expectedSeconds: 0,
		},
		{
			name:            "entirely after range",
			begin:           "2026-07-16T18:00:00Z",
			end:             "2026-07-16T19:00:00Z",
			expectedSeconds: 0,
		},
		{
			name:            "ends at range start",
			begin:           "2026-07-16T08:00:00Z",
			end:             "2026-07-16T09:00:00Z",
			expectedSeconds: 0,
		},
		{
			name:            "begins at range end",
			begin:           "2026-07-16T17:00:00Z",
			end:             "2026-07-16T18:00:00Z",
			expectedSeconds: 0,
		},
		{
			name:            "zero length interval",
			begin:           "2026-07-16T12:00:00Z",
			end:             "2026-07-16T12:00:00Z",
			expectedSeconds: 0,
		},
		{
			name:            "reversed interval",
			begin:           "2026-07-16T13:00:00Z",
			end:             "2026-07-16T12:00:00Z",
			expectedSeconds: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			begin := timestamp(t, tt.begin)
			end := timestamp(t, tt.end)

			// WHEN
			got := overlappingSeconds(begin, end, rangeStart, rangeEnd)

			// THEN
			assert.Equal(t, tt.expectedSeconds, got)
		})
	}
}

func TestSecondsTrackedToday(t *testing.T) {
	now := timestamp(t, "2026-01-16T14:00:00+01:00")
	taskLogsWithinDay := []TaskLogEntry{
		{
			BeginTS: timestamp(t, "2026-01-16T09:15:00+01:00"),
			EndTS:   timestamp(t, "2026-01-16T10:45:00+01:00"),
		},
		{
			BeginTS: timestamp(t, "2026-01-16T11:30:00+01:00"),
			EndTS:   timestamp(t, "2026-01-16T12:15:00+01:00"),
		},
	}
	taskLogsCrossingMidnight := []TaskLogEntry{
		{
			BeginTS: timestamp(t, "2026-01-15T23:15:00+01:00"),
			EndTS:   timestamp(t, "2026-01-16T00:45:00+01:00"),
		},
	}
	zeroTimestamp := time.Time{}
	activeTaskLogBeginTSWithinDay := timestamp(t, "2026-01-16T13:00:00+01:00")
	activeTaskLogBeginTSBeforeMidnight := timestamp(t, "2026-01-15T23:15:00+01:00")

	testCases := []struct {
		name                 string
		finishedTaskLogs     []TaskLogEntry
		activeTaskLogBeginTS *time.Time
		expectedSeconds      int
	}{
		{
			name:            "no task logs",
			expectedSeconds: 0,
		},
		{
			name:                 "zero active task log begin timestamp",
			activeTaskLogBeginTS: &zeroTimestamp,
			expectedSeconds:      0,
		},
		{
			name:             "task logs within today",
			finishedTaskLogs: taskLogsWithinDay,
			expectedSeconds:  (90 + 45) * 60,
		},
		{
			name:                 "task logs within today and an active task",
			finishedTaskLogs:     taskLogsWithinDay,
			activeTaskLogBeginTS: &activeTaskLogBeginTSWithinDay,
			expectedSeconds:      (90 + 45 + 60) * 60,
		},
		{
			name:             "task log crossing midnight",
			finishedTaskLogs: taskLogsCrossingMidnight,
			expectedSeconds:  45 * 60,
		},
		{
			name:                 "active task log beginning before midnight",
			activeTaskLogBeginTS: &activeTaskLogBeginTSBeforeMidnight,
			expectedSeconds:      14 * 60 * 60,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			// WHEN
			got := SecondsTrackedToday(
				tt.finishedTaskLogs,
				tt.activeTaskLogBeginTS,
				now,
			)

			// THEN
			assert.Equal(t, tt.expectedSeconds, got)
		})
	}
}

func timestamp(t *testing.T, value string) time.Time {
	t.Helper()

	timestamp, err := time.Parse(time.RFC3339, value)
	require.NoError(t, err)

	return timestamp
}
