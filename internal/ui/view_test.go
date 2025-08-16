package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDurationValidityContext(t *testing.T) {
	testCases := []struct {
		name             string
		beginTS          string
		endTS            string
		expectedCtx      string
		expectedValidity tlFormValidity
	}{
		// success cases
		{
			name:             "less than an hour",
			beginTS:          "2025/08/08 00:40",
			endTS:            "2025/08/08 00:48",
			expectedCtx:      "You're recording 8m",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "exact hour",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 01:00",
			expectedCtx:      "You're recording 1h",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "hours and minutes",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 02:30",
			expectedCtx:      "You're recording 2h 30m",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "across day boundary",
			beginTS:          "2025/08/08 23:30",
			endTS:            "2025/08/09 00:15",
			expectedCtx:      "You're recording 45m",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "exactly at 8h threshold",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 08:00",
			expectedCtx:      "You're recording 8h",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "> 8h threshold",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 08:01",
			expectedCtx:      "You're recording 8h 1m",
			expectedValidity: tlSubmitWarn,
		},
		{
			name:             "very long duration",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/09 02:00",
			expectedCtx:      "You're recording 26h",
			expectedValidity: tlSubmitWarn,
		},
		// failure cases
		{
			name:             "empty begin time",
			beginTS:          "",
			endTS:            "2025/08/08 00:10",
			expectedCtx:      "Begin time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "empty end time",
			beginTS:          "2025/08/08 00:10",
			endTS:            "",
			expectedCtx:      "End time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "begin time as whitespace only",
			beginTS:          "   ",
			endTS:            "2025/08/08 00:10",
			expectedCtx:      "Begin time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "end time as whitespace only",
			beginTS:          "2025/08/08 00:10",
			endTS:            "   ",
			expectedCtx:      "End time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "invalid begin ts",
			beginTS:          "2025-08-08 00:10",
			endTS:            "2025/08/08 00:20",
			expectedCtx:      "Begin time is invalid",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "invalid end format",
			beginTS:          "2025/08/08 00:10",
			endTS:            "08-08-2025 00:20",
			expectedCtx:      "End time is invalid",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "end before start",
			beginTS:          "2025/08/08 01:00",
			endTS:            "2025/08/08 00:59",
			expectedCtx:      "End time is before begin time",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "zero duration",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 00:00",
			expectedCtx:      "You're recording no time, change begin and/or end time",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "one minute duration",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 00:01",
			expectedCtx:      "You're recording 1m",
			expectedValidity: tlSubmitOk,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gotCtx, gotValidity := getDurationValidityContext(tt.beginTS, tt.endTS)

			assert.Equal(t, tt.expectedCtx, gotCtx)
			assert.Equal(t, tt.expectedValidity, gotValidity)
		})
	}
}
