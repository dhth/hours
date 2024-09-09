package persistence

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/dhth/hours/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // sqlite driver
)

func TestRepository(t *testing.T) {
	testDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("error opening DB: %v", err)
	}

	err = InitDB(testDB)
	if err != nil {
		t.Fatalf("error initializing DB: %v", err)
	}
	err = UpgradeDB(testDB, 1)
	if err != nil {
		t.Fatalf("error upgrading DB: %v", err)
	}

	t.Run("TestInsertTask", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		seedDB(t, testDB)

		// WHEN
		summary := "task 1"
		taskID, err := InsertTask(testDB, summary)

		// THEN
		require.NoError(t, err, "failed to insert task")

		task, fetchErr := fetchTaskByID(testDB, int(taskID))
		require.NoError(t, fetchErr, "failed to fetch task")

		assert.Equal(t, 3, task.ID)
		assert.Equal(t, summary, task.Summary)
		assert.True(t, task.Active)
		assert.Equal(t, 0, task.SecsSpent)
	})

	t.Run("TestUpdateActiveTL", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		seedDB(t, testDB)
		taskID := 1
		numSeconds := 60 * 90
		endTS := time.Now()
		beginTS := endTS.Add(time.Second * -1 * time.Duration(numSeconds))
		tlID, insertErr := InsertNewTL(testDB, taskID, beginTS)
		require.NoError(t, insertErr, "failed to insert task log")

		taskBefore, err := fetchTaskByID(testDB, taskID)
		require.NoError(t, err, "failed to fetch task")
		numSecondsBefore := taskBefore.SecsSpent

		// WHEN
		comment := "a task log"
		err = UpdateActiveTL(testDB, int(tlID), taskID, beginTS, endTS, numSeconds, comment)

		// THEN
		require.NoError(t, err, "failed to update task log")

		taskLog, err := fetchTaskLogByID(testDB, int(tlID))
		require.NoError(t, err, "failed to fetch task log")

		taskAfter, err := fetchTaskByID(testDB, taskID)
		require.NoError(t, err, "failed to fetch task")

		assert.Equal(t, numSeconds, taskLog.SecsSpent)
		assert.Equal(t, comment, taskLog.Comment)
		assert.Equal(t, numSecondsBefore+numSeconds, taskAfter.SecsSpent)
	})

	t.Run("TestInsertManualTL", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		seedDB(t, testDB)
		taskID := 1

		taskBefore, err := fetchTaskByID(testDB, taskID)
		require.NoError(t, err, "failed to fetch task")
		numSecondsBefore := taskBefore.SecsSpent

		// WHEN
		comment := "a task log"
		numSeconds := 60 * 90
		endTS := time.Now()
		beginTS := endTS.Add(time.Second * -1 * time.Duration(numSeconds))
		tlID, err := InsertManualTL(testDB, taskID, beginTS, endTS, comment)

		// THEN
		require.NoError(t, err, "failed to insert task log")

		taskLog, err := fetchTaskLogByID(testDB, int(tlID))
		require.NoError(t, err, "failed to fetch task log")

		taskAfter, err := fetchTaskByID(testDB, taskID)
		require.NoError(t, err, "failed to fetch task")

		assert.Equal(t, numSeconds, taskLog.SecsSpent)
		assert.Equal(t, comment, taskLog.Comment)
		assert.Equal(t, numSecondsBefore+numSeconds, taskAfter.SecsSpent)
	})

	t.Run("TestDeleteTaskLogEntry", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		seedDB(t, testDB)
		taskID := 1
		tlID := 1
		taskBefore, err := fetchTaskByID(testDB, taskID)
		require.NoError(t, err, "failed to fetch task")
		numSecondsBefore := taskBefore.SecsSpent
		taskLog, err := fetchTaskLogByID(testDB, tlID)
		require.NoError(t, err, "failed to fetch task log")

		// WHEN
		err = DeleteTaskLogEntry(testDB, &taskLog)

		// THEN
		require.NoError(t, err, "failed to insert task log")

		taskAfter, err := fetchTaskByID(testDB, taskID)
		require.NoError(t, err, "failed to fetch task")

		assert.Equal(t, numSecondsBefore-taskLog.SecsSpent, taskAfter.SecsSpent)
	})

	err = testDB.Close()
	if err != nil {
		t.Fatalf("error closing DB: %v", err)
	}
}

func cleanupDB(t *testing.T, testDB *sql.DB) {
	var err error
	for _, tbl := range []string{"task_log", "task"} {
		_, err = testDB.Exec(fmt.Sprintf("DELETE FROM %s", tbl))
		if err != nil {
			t.Fatalf("failed to clean up table %q: %v", tbl, err)
		}
		_, err := testDB.Exec("DELETE FROM sqlite_sequence WHERE name=?;", tbl)
		if err != nil {
			t.Fatalf("failed to reset auto increment for table %q: %v", tbl, err)
		}
	}
}

func seedDB(t *testing.T, db *sql.DB) {
	t.Helper()

	ua := time.Now().UTC()
	ca := ua.Add(time.Hour * 24 * 7 * -1)
	hour := 60 * 60
	tasks := []types.Task{
		{
			ID:        1,
			Summary:   "seeded task 1",
			Active:    true,
			CreatedAt: ca,
			UpdatedAt: ua,
			SecsSpent: 5 * hour,
		},
		{
			ID:        2,
			Summary:   "seeded task 2",
			Active:    true,
			CreatedAt: ca,
			UpdatedAt: ua,
			SecsSpent: 4 * hour,
		},
	}

	taskLogs := []types.TaskLogEntry{
		{
			ID:        1,
			TaskID:    1,
			BeginTS:   ca.Add(time.Hour * 2),
			EndTS:     ca.Add(time.Hour * 4),
			SecsSpent: 2 * hour,
			Comment:   "task 1 tl 1",
		},
		{
			ID:        2,
			TaskID:    1,
			BeginTS:   ca.Add(time.Hour * 6),
			EndTS:     ca.Add(time.Hour * 9),
			SecsSpent: 3 * hour,
			Comment:   "task 1 tl 2",
		},
		{
			ID:        3,
			TaskID:    2,
			BeginTS:   ca.Add(time.Hour * 2),
			EndTS:     ca.Add(time.Hour * 6),
			SecsSpent: 4 * hour,
			Comment:   "task 2 tl 1",
		},
	}

	for _, task := range tasks {
		_, err := db.Exec(`
INSERT INTO task (id, summary, secs_spent, active, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, task.ID, task.Summary, task.SecsSpent, task.Active, task.CreatedAt, task.UpdatedAt)
		if err != nil {
			t.Fatalf("failed to insert data into table \"task\": %v", err)
		}
	}

	for _, taskLog := range taskLogs {
		_, err := db.Exec(`
INSERT INTO task_log (id, task_id, begin_ts, end_ts, secs_spent, comment, active)
VALUES (?, ?, ?, ?, ?, ?, ?)`, taskLog.ID, taskLog.TaskID, taskLog.BeginTS, taskLog.EndTS, taskLog.SecsSpent, taskLog.Comment, false)
		if err != nil {
			t.Fatalf("failed to insert data into table \"task_log\": %v", err)
		}
	}
}
