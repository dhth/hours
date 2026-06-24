package ui

import (
	"database/sql"
	"testing"
	"time"

	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/ui/theme"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // sqlite driver
)

func TestGetTaskLogNoTruncate(t *testing.T) {
	// GIVEN
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	require.NoError(t, pers.InitDB(db))

	longTaskSummary := "this-is-a-very-long-task-summary"
	longComment := "this is a very long comment that should not be truncated"
	taskID, err := pers.InsertTask(db, longTaskSummary)
	require.NoError(t, err)

	beginTS := time.Date(2026, time.January, 2, 10, 0, 0, 0, time.UTC)
	endTS := beginTS.Add(90 * time.Minute)
	_, err = pers.InsertManualTL(db, taskID, beginTS, endTS, &longComment)
	require.NoError(t, err)

	style := NewStyle(theme.Default())
	start := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	// WHEN
	truncatedLog, err := getTaskLog(db, style, start, end, types.TaskStatusAny, logLimit, true, false)
	require.NoError(t, err)
	untruncatedLog, err := getTaskLog(db, style, start, end, types.TaskStatusAny, logLimit, true, true)
	require.NoError(t, err)

	// THEN
	require.NotContains(t, truncatedLog, longTaskSummary)
	require.NotContains(t, truncatedLog, longComment)
	require.Contains(t, truncatedLog, longTaskSummary[:logTaskCharsBudget])
	require.Contains(t, untruncatedLog, longTaskSummary)
	require.Contains(t, untruncatedLog, longComment)
}
