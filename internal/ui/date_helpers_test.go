package ui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateDuration(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedStartStr string
		expectedEndStr   string
		expectedNumDays  int
		ok               bool
	}{
		// success
		{
			name:             "a range of 1 day",
			input:            "2024/06/10...2024/06/11",
			expectedStartStr: "2024/06/10 00:00",
			expectedEndStr:   "2024/06/11 00:00",
			expectedNumDays:  2,
			ok:               true,
		},
		{
			name:             "a range of 2 days",
			input:            "2024/06/29...2024/07/01",
			expectedStartStr: "2024/06/29 00:00",
			expectedEndStr:   "2024/07/01 00:00",
			expectedNumDays:  3,
			ok:               true,
		},
		{
			name:             "a range of 1 year",
			input:            "2024/06/29...2025/06/29",
			expectedStartStr: "2024/06/29 00:00",
			expectedEndStr:   "2025/06/29 00:00",
			expectedNumDays:  366,
			ok:               true,
		},
		// failures
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "only one date",
			input: "2024/06/10",
		},
		{
			name:  "badly formatted start date",
			input: "2024/0610...2024/06/10",
		},
		{
			name:  "badly formatted end date",
			input: "2024/06/10...2024/0610",
		},
		{
			name:  "a range of 0 days",
			input: "2024/06/10...2024/06/10",
		},
		{
			name:  "end date before start date",
			input: "2024/06/10...2024/06/08",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseDateDuration(tt.input)

			if tt.ok {
				startStr := got.start.Format(timeFormat)
				endStr := got.end.Format(timeFormat)

				assert.True(t, ok)
				assert.Equal(t, tt.expectedStartStr, startStr)
				assert.Equal(t, tt.expectedEndStr, endStr)
				assert.Equal(t, tt.expectedNumDays, got.numDays)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestGetTimePeriod(t *testing.T) {
	now, err := time.ParseInLocation(string(timeFormat), "2024/06/20 20:00", time.Local)
	if err != nil {
		t.Fatalf("error setting up the test: time is not valid: %s", err)
	}

	nowME, err := time.ParseInLocation(string(timeFormat), "2024/05/31 20:00", time.Local)
	if err != nil {
		t.Fatalf("error setting up the test: time is not valid: %s", err)
	}

	nowMB, err := time.ParseInLocation(string(timeFormat), "2024/06/01 20:00", time.Local)
	if err != nil {
		t.Fatalf("error setting up the test: time is not valid: %s", err)
	}

	testCases := []struct {
		name             string
		period           string
		now              time.Time
		fullWeek         bool
		expectedStartStr string
		expectedEndStr   string
		expectedNumDays  int
		err              error
	}{
		// success
		{
			name:             "today",
			period:           "today",
			now:              now,
			expectedStartStr: "2024/06/20 00:00",
			expectedEndStr:   "2024/06/21 00:00",
			expectedNumDays:  1,
		},
		{
			name:             "'today' at end of month",
			period:           "today",
			now:              nowME,
			expectedStartStr: "2024/05/31 00:00",
			expectedEndStr:   "2024/06/01 00:00",
			expectedNumDays:  1,
		},
		{
			name:             "'yest' at beginning of month",
			period:           "yest",
			now:              nowMB,
			expectedStartStr: "2024/05/31 00:00",
			expectedEndStr:   "2024/06/01 00:00",
			expectedNumDays:  1,
		},
		{
			name:             "3d",
			period:           "3d",
			now:              now,
			expectedStartStr: "2024/06/18 00:00",
			expectedEndStr:   "2024/06/21 00:00",
			expectedNumDays:  3,
		},
		{
			name:             "week",
			period:           "week",
			now:              now,
			expectedStartStr: "2024/06/17 00:00",
			expectedEndStr:   "2024/06/21 00:00",
			expectedNumDays:  4,
		},
		{
			name:             "full week",
			period:           "week",
			now:              now,
			fullWeek:         true,
			expectedStartStr: "2024/06/17 00:00",
			expectedEndStr:   "2024/06/24 00:00",
			expectedNumDays:  7,
		},
		{
			name:             "a date",
			period:           "2024/06/20",
			expectedStartStr: "2024/06/20 00:00",
			expectedEndStr:   "2024/06/21 00:00",
			expectedNumDays:  1,
		},
		{
			name:             "a date range",
			period:           "2024/06/15...2024/06/20",
			expectedStartStr: "2024/06/15 00:00",
			expectedEndStr:   "2024/06/21 00:00",
			expectedNumDays:  6,
		},
		// failures
		{
			name:   "a faulty date",
			period: "2024/06-15",
			err:    timePeriodNotValidErr,
		},
		{
			name:   "a faulty date range",
			period: "2024/06/15...2024",
			err:    timePeriodNotValidErr,
		},
		{
			name:   "a date range too large",
			period: "2024/06/15...2024/06/22",
			err:    timePeriodTooLargeErr,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTimePeriod(tt.period, tt.now, tt.fullWeek)

			startStr := got.start.Format(timeFormat)
			endStr := got.end.Format(timeFormat)

			if tt.err == nil {
				assert.Equal(t, tt.expectedStartStr, startStr)
				assert.Equal(t, tt.expectedEndStr, endStr)
				assert.Equal(t, tt.expectedNumDays, got.numDays)
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.err, err)
			}
		})
	}
}

func TestGetTSRelative(t *testing.T) {
	reference := time.Date(2024, 6, 29, 12, 0, 0, 0, time.Local)
	testCases := []struct {
		name      string
		ts        time.Time
		reference time.Time
		expected  tsRelative
	}{
		{
			name:      "ts in the future",
			ts:        time.Date(2024, 6, 30, 6, 0, 0, 0, time.Local),
			reference: reference,
			expected:  tsFromFuture,
		},
		{
			name:      "ts on the same day as the reference",
			ts:        time.Date(2024, 6, 29, 6, 0, 0, 0, time.Local),
			reference: reference,
			expected:  tsFromToday,
		},
		{
			name:      "ts from a day before the reference",
			ts:        time.Date(2024, 6, 28, 23, 59, 0, 0, time.Local),
			reference: reference,
			expected:  tsFromYesterday,
		},
		{
			name:      "ts from the first day of the week",
			ts:        time.Date(2024, 6, 24, 0, 1, 0, 0, time.Local),
			reference: reference,
			expected:  tsFromThisWeek,
		},
		{
			name:      "ts from before the week",
			ts:        time.Date(2024, 6, 23, 23, 59, 0, 0, time.Local),
			reference: reference,
			expected:  tsFromBeforeThisWeek,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := getTSRelative(tt.ts, tt.reference)
			assert.Equal(t, tt.expected, got)
		})
	}
}
