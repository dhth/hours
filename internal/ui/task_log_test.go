package ui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTSRelative(t *testing.T) {
	// GIVEN
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
			// WHEN
			got := getTSRelative(tt.ts, tt.reference)

			// THEN
			assert.Equal(t, tt.expected, got)
		})
	}
}
