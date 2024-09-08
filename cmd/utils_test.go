package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandTilde(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		homeDir  string
		expected string
	}{
		{
			name:     "a simple case",
			path:     "~/some/path",
			homeDir:  "/Users/trinity",
			expected: "/Users/trinity/some/path",
		},
		{
			name:     "path with no ~",
			path:     "some/path",
			homeDir:  "/Users/trinity",
			expected: "some/path",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := expandTilde(tt.path, tt.homeDir)

			assert.Equal(t, tt.expected, got)
		})
	}
}
