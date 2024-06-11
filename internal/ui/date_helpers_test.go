package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDateDuration(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedStartStr string
		expectedEndStr   string
		err              error
	}{
		// success
		{
			name:             "a range of 1 day",
			input:            "2024/06/10...2024/06/11",
			expectedStartStr: "2024/06/10 00:00",
			expectedEndStr:   "2024/06/11 00:00",
		},
		{
			name:             "a range of 3 days",
			input:            "2024/06/29...2024/07/01",
			expectedStartStr: "2024/06/29 00:00",
			expectedEndStr:   "2024/07/01 00:00",
		},
		// failures
		{
			name:  "empty string",
			input: "",
			err:   timePeriodNotValidErr,
		},
		{
			name:  "only one date",
			input: "2024/06/10",
			err:   timePeriodNotValidErr,
		},
		{
			name:  "badly formatted start date",
			input: "2024/0610...2024/06/10",
			err:   timePeriodNotValidErr,
		},
		{
			name:  "badly formatted end date",
			input: "2024/06/10...2024/0610",
			err:   timePeriodNotValidErr,
		},
		{
			name:  "a range of 0 days",
			input: "2024/06/10...2024/06/10",
			err:   timePeriodNotValidErr,
		},
		{
			name:  "end date before start date",
			input: "2024/06/10...2024/06/08",
			err:   timePeriodNotValidErr,
		},
		{
			name:  "a range of 8 days",
			input: "2024/06/29...2024/07/06",
			err:   timePeriodTooLargeErr,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateDuration(tt.input)

			startStr := got.start.Format(timeFormat)
			endStr := got.end.Format(timeFormat)

			if tt.err == nil {
				assert.Equal(t, tt.expectedStartStr, startStr)
				assert.Equal(t, tt.expectedEndStr, endStr)
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.err, err)
			}

		})
	}

}
