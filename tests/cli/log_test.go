package cli

import (
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	fx := NewFixture(t, testBinaryPath)
	now := time.Date(2025, time.October, 24, 12, 0, 0, 0, time.UTC)

	_, err := fx.RunGen(42, now)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		period string
	}{
		{name: "today", period: "today"},
		{name: "yest", period: "yest"},
		{name: "3d", period: "3d"},
		{name: "week", period: "week"},
		{name: "date", period: "2025/10/24"},
		{name: "date range", period: "2025/10/20...2025/10/24"},
		{name: "incorrect argument", period: "blah"},
		{name: "incorrect date", period: "2025/1024"},
		{name: "incorrect date range", period: "2025/1024...2025/10/24"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCmd([]string{"log", "--plain"})
			cmd.AddArgs(tc.period)
			cmd.SetEnv("HOURS_NOW", now.Format(time.RFC3339))
			cmd.UseDB()

			result, runErr := fx.RunCmd(cmd)

			require.NoError(t, runErr)
			snaps.MatchStandaloneSnapshot(t, result)
		})
	}
}
