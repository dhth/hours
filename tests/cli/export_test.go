package cli

import (
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestReportExport(t *testing.T) {
	// GIVEN
	fx := NewFixture(t, testBinaryPath)
	now := time.Date(2025, time.October, 24, 12, 0, 0, 0, time.UTC)

	_, err := fx.RunGen(42, now)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		args   []string
		period string
	}{
		{name: "json today", args: []string{"report", "--format", "json"}, period: "today"},
		{name: "csv today", args: []string{"report", "--format", "csv"}, period: "today"},
		{name: "json 3d", args: []string{"report", "--format", "json"}, period: "3d"},
		{name: "csv 3d", args: []string{"report", "--format", "csv"}, period: "3d"},
		{name: "json agg today", args: []string{"report", "--format", "json", "--agg"}, period: "today"},
		{name: "csv agg today", args: []string{"report", "--format", "csv", "--agg"}, period: "today"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN
			cmd := NewCmd(tc.args)
			cmd.AddArgs(tc.period)
			cmd.SetEnv("HOURS_NOW", now.Format(time.RFC3339))
			cmd.UseDB()

			result, runErr := fx.RunCmd(cmd)

			// THEN
			require.NoError(t, runErr)
			snaps.MatchStandaloneSnapshot(t, result)
		})
	}
}

func TestLogExport(t *testing.T) {
	// GIVEN
	fx := NewFixture(t, testBinaryPath)
	now := time.Date(2025, time.October, 24, 12, 0, 0, 0, time.UTC)

	_, err := fx.RunGen(42, now)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		format string
		period string
	}{
		{name: "json today", format: "json", period: "today"},
		{name: "csv today", format: "csv", period: "today"},
		{name: "json 3d", format: "json", period: "3d"},
		{name: "csv 3d", format: "csv", period: "3d"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN
			cmd := NewCmd([]string{"log", "--format", tc.format})
			cmd.AddArgs(tc.period)
			cmd.SetEnv("HOURS_NOW", now.Format(time.RFC3339))
			cmd.UseDB()

			result, runErr := fx.RunCmd(cmd)

			// THEN
			require.NoError(t, runErr)
			snaps.MatchStandaloneSnapshot(t, result)
		})
	}
}

func TestStatsExport(t *testing.T) {
	// GIVEN
	fx := NewFixture(t, testBinaryPath)
	now := time.Date(2025, time.October, 24, 12, 0, 0, 0, time.UTC)

	_, err := fx.RunGen(42, now)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		format string
		period string
	}{
		{name: "json today", format: "json", period: "today"},
		{name: "csv today", format: "csv", period: "today"},
		{name: "json 3d", format: "json", period: "3d"},
		{name: "csv 3d", format: "csv", period: "3d"},
		{name: "json all", format: "json", period: "all"},
		{name: "csv all", format: "csv", period: "all"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN
			cmd := NewCmd([]string{"stats", "--format", tc.format})
			cmd.AddArgs(tc.period)
			cmd.SetEnv("HOURS_NOW", now.Format(time.RFC3339))
			cmd.UseDB()

			result, runErr := fx.RunCmd(cmd)

			// THEN
			require.NoError(t, runErr)
			snaps.MatchStandaloneSnapshot(t, result)
		})
	}
}

func TestExportFormatErrors(t *testing.T) {
	// GIVEN
	fx := NewFixture(t, testBinaryPath)
	now := time.Date(2025, time.October, 24, 12, 0, 0, 0, time.UTC)

	_, err := fx.RunGen(42, now)
	require.NoError(t, err)

	testCases := []struct {
		name string
		args []string
	}{
		{name: "report interactive json", args: []string{"report", "--interactive", "--format", "json", "today"}},
		{name: "log interactive json", args: []string{"log", "--interactive", "--format", "json", "today"}},
		{name: "stats interactive json", args: []string{"stats", "--interactive", "--format", "json", "today"}},
		{name: "report incorrect format", args: []string{"report", "--format", "xml", "today"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// WHEN
			cmd := NewCmd(tc.args)
			cmd.SetEnv("HOURS_NOW", now.Format(time.RFC3339))
			cmd.UseDB()

			result, runErr := fx.RunCmd(cmd)

			// THEN
			require.NoError(t, runErr)
			snaps.MatchStandaloneSnapshot(t, result)
		})
	}
}
