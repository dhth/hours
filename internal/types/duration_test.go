package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTaskLogTimes(t *testing.T) {
	testCases := []struct {
		name     string
		beginStr string
		endStr   string
		err      error
	}{
		// Successes
		{
			name:     "valid times - less than an hour",
			beginStr: "2025/08/08 00:40",
			endStr:   "2025/08/08 00:48",
		},
		{
			name:     "valid times - exact hour",
			beginStr: "2025/08/08 00:00",
			endStr:   "2025/08/08 01:00",
		},
		{
			name:     "valid times - hours and minutes",
			beginStr: "2025/08/08 00:00",
			endStr:   "2025/08/08 02:30",
		},
		{
			name:     "valid times - across day boundary",
			beginStr: "2025/08/08 23:30",
			endStr:   "2025/08/09 00:15",
		},
		{
			name:     "valid times - exactly at 8h",
			beginStr: "2025/08/08 00:00",
			endStr:   "2025/08/08 08:00",
		},
		{
			name:     "valid times - very long duration",
			beginStr: "2025/08/08 00:00",
			endStr:   "2025/08/09 02:00",
		},
		{
			name:     "valid times - exactly one minute",
			beginStr: "2025/08/08 00:00",
			endStr:   "2025/08/08 00:01",
		},
		// Failures
		{
			name:     "empty begin time",
			beginStr: "",
			endStr:   "2025/08/08 00:10",
			err:      errBeginTimeIsEmpty,
		},
		{
			name:     "empty end time",
			beginStr: "2025/08/08 00:10",
			endStr:   "",
			err:      errEndTimeIsEmpty,
		},
		{
			name:     "begin time as whitespace only",
			beginStr: "   ",
			endStr:   "2025/08/08 00:10",
			err:      errBeginTimeIsEmpty,
		},
		{
			name:     "end time as whitespace only",
			beginStr: "2025/08/08 00:10",
			endStr:   "   ",
			err:      errEndTimeIsEmpty,
		},
		{
			name:     "invalid begin time format",
			beginStr: "2025-08-08 00:10",
			endStr:   "2025/08/08 00:20",
			err:      errBeginTimeIsInvalid,
		},
		{
			name:     "invalid end time format",
			beginStr: "2025/08/08 00:10",
			endStr:   "08-08-2025 00:20",
			err:      errEndTimeIsInvalid,
		},
		{
			name:     "end time before begin time",
			beginStr: "2025/08/08 01:00",
			endStr:   "2025/08/08 00:59",
			err:      errEndTimeBeforeBeginTime,
		},
		{
			name:     "zero duration",
			beginStr: "2025/08/08 00:00",
			endStr:   "2025/08/08 00:00",
			err:      ErrDurationNotLongEnough,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			beginTS, endTS, err := ParseTaskLogTimes(tt.beginStr, tt.endStr)

			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.False(t, beginTS.IsZero())
				assert.False(t, endTS.IsZero())
			}
		})
	}
}
