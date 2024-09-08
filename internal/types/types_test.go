package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeDuration(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "0 seconds",
			input:    0,
			expected: "0s",
		},
		{
			name:     "30 seconds",
			input:    30,
			expected: "30s",
		},
		{
			name:     "60 seconds",
			input:    60,
			expected: "1m",
		},
		{
			name:     "1805 seconds",
			input:    1805,
			expected: "30m",
		},
		{
			name:     "3605 seconds",
			input:    3605,
			expected: "1h",
		},
		{
			name:     "4200 seconds",
			input:    4200,
			expected: "1h 10m",
		},
		{
			name:     "87000 seconds",
			input:    87000,
			expected: "24h 10m",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := HumanizeDuration(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
