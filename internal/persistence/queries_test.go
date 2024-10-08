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

const (
	secsInOneHour  = 60 * 60
	taskLogComment = "a task log outside the time range"
)

func TestRepository(t *testing.T) {
	testDB, err := sql.Open("sqlite", ":memory:")
	require.NoErrorf(t, err, "error opening DB: %v", err)

	err = InitDB(testDB)
	require.NoErrorf(t, err, "error initializing DB: %v", err)

	err = UpgradeDB(testDB, 1)
	require.NoErrorf(t, err, "error upgrading DB: %v", err)

	t.Run("TestInsertTask", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		referenceTS := time.Now()
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)

		// WHEN
		summary := "task 1"
		taskID, err := InsertTask(testDB, summary)

		// THEN
		require.NoError(t, err, "failed to insert task")

		task, fetchErr := fetchTaskByID(testDB, taskID)
		require.NoError(t, fetchErr, "failed to fetch task")

		assert.Equal(t, 3, task.ID)
		assert.Equal(t, summary, task.Summary)
		assert.True(t, task.Active)
		assert.Zero(t, task.SecsSpent)
	})

	t.Run("TestUpdateActiveTL", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		referenceTS := time.Now()
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)
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
		err = UpdateActiveTL(testDB, tlID, taskID, beginTS, endTS, numSeconds, comment)

		// THEN
		require.NoError(t, err, "failed to update task log")

		taskLog, err := fetchTaskLogByID(testDB, tlID)
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
		referenceTS := time.Now()
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)
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

		taskLog, err := fetchTaskLogByID(testDB, tlID)
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
		referenceTS := time.Now()
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)
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

	t.Run("TestFetchTLEntriesBetweenTS", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		referenceTS := time.Date(2024, time.September, 1, 9, 0, 0, 0, time.Local)
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)

		taskID := 1
		numSeconds := 60 * 90
		tlEndTS := referenceTS.Add(time.Hour * 2)
		tlBeginTS := tlEndTS.Add(time.Second * -1 * time.Duration(numSeconds))
		_, err = InsertManualTL(testDB, taskID, tlBeginTS, tlEndTS, taskLogComment)
		require.NoError(t, err, "failed to insert task log")

		// WHEN
		reportBeginTS := referenceTS.Add(time.Hour * 24 * 7 * -2)
		entries, err := FetchTLEntriesBetweenTS(testDB, reportBeginTS, referenceTS, 100)

		// THEN
		require.NoError(t, err, "failed to fetch report entries")
		require.Len(t, entries, 3)
	})

	t.Run("TestFetchStats", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		referenceTS := time.Date(2024, time.September, 1, 9, 0, 0, 0, time.Local)
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)

		taskID := 1
		comment := "an extra task log"
		numSeconds := 60 * 90
		tlEndTS := referenceTS.Add(time.Hour * 2)
		tlBeginTS := tlEndTS.Add(time.Second * -1 * time.Duration(numSeconds))
		_, err = InsertManualTL(testDB, taskID, tlBeginTS, tlEndTS, comment)
		require.NoError(t, err, "failed to insert task log")

		// WHEN
		entries, err := FetchStats(testDB, 100)

		// THEN
		require.NoError(t, err, "failed to fetch report entries")
		require.Len(t, entries, 2)

		assert.Equal(t, 1, entries[0].TaskID)
		assert.Equal(t, 3, entries[0].NumEntries)
		assert.Equal(t, 5*secsInOneHour+numSeconds, entries[0].SecsSpent)

		assert.Equal(t, 2, entries[1].TaskID)
		assert.Equal(t, 1, entries[1].NumEntries)
		assert.Equal(t, 4*secsInOneHour, entries[1].SecsSpent)
	})

	t.Run("TestFetchStatsBetweenTS", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		referenceTS := time.Date(2024, time.September, 1, 9, 0, 0, 0, time.Local)
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)

		taskID := 1
		numSeconds := 60 * 90
		tlEndTS := referenceTS.Add(time.Hour * 2)
		tlBeginTS := tlEndTS.Add(time.Second * -1 * time.Duration(numSeconds))
		_, err = InsertManualTL(testDB, taskID, tlBeginTS, tlEndTS, taskLogComment)
		require.NoError(t, err, "failed to insert task log")

		// WHEN
		reportBeginTS := referenceTS.Add(time.Hour * 24 * 7 * -2)
		entries, err := FetchStatsBetweenTS(testDB, reportBeginTS, referenceTS, 100)

		// THEN
		require.NoError(t, err, "failed to fetch report entries")
		require.Len(t, entries, 2)

		assert.Equal(t, 1, entries[0].TaskID)
		assert.Equal(t, 2, entries[0].NumEntries)
		assert.Equal(t, 5*secsInOneHour, entries[0].SecsSpent)

		assert.Equal(t, 2, entries[1].TaskID)
		assert.Equal(t, 1, entries[1].NumEntries)
		assert.Equal(t, 4*secsInOneHour, entries[1].SecsSpent)
	})

	t.Run("TestFetchReportBetweenTS", func(t *testing.T) {
		t.Cleanup(func() { cleanupDB(t, testDB) })

		// GIVEN
		referenceTS := time.Date(2024, time.September, 1, 9, 0, 0, 0, time.Local)
		seedData := getTestData(referenceTS)
		seedDB(t, testDB, seedData)

		taskID := 1
		numSeconds := 60 * 90
		tlEndTS := referenceTS.Add(time.Hour * 2)
		tlBeginTS := tlEndTS.Add(time.Second * -1 * time.Duration(numSeconds))
		_, err = InsertManualTL(testDB, taskID, tlBeginTS, tlEndTS, taskLogComment)
		require.NoError(t, err, "failed to insert task log")

		// WHEN
		reportBeginTS := referenceTS.Add(time.Hour * 24 * 7 * -2)
		entries, err := FetchReportBetweenTS(testDB, reportBeginTS, referenceTS, 100)

		// THEN
		require.NoError(t, err, "failed to fetch report entries")

		require.Len(t, entries, 2)
		assert.Equal(t, 2, entries[0].TaskID)
		assert.Equal(t, 1, entries[0].NumEntries)
		assert.Equal(t, 4*secsInOneHour, entries[0].SecsSpent)

		assert.Equal(t, 1, entries[1].TaskID)
		assert.Equal(t, 2, entries[1].NumEntries)
		assert.Equal(t, 5*secsInOneHour, entries[1].SecsSpent)
	})

	err = testDB.Close()
	require.NoErrorf(t, err, "error closing DB: %v", err)
}

func cleanupDB(t *testing.T, testDB *sql.DB) {
	t.Helper()

	var err error
	for _, tbl := range []string{"task_log", "task"} {
		_, err = testDB.Exec(fmt.Sprintf("DELETE FROM %s", tbl))
		require.NoErrorf(t, err, "failed to clean up table %q: %v", tbl, err)

		_, err := testDB.Exec("DELETE FROM sqlite_sequence WHERE name=?;", tbl)
		require.NoErrorf(t, err, "failed to reset auto increment for table %q: %v", tbl, err)
	}
}

type testData struct {
	tasks    []types.Task
	taskLogs []types.TaskLogEntry
}

func getTestData(referenceTS time.Time) testData {
	ua := referenceTS.UTC()
	ca := ua.Add(time.Hour * 24 * 7 * -1)
	tasks := []types.Task{
		{
			ID:        1,
			Summary:   "seeded task 1",
			Active:    true,
			CreatedAt: ca,
			UpdatedAt: ca.Add(time.Hour * 9),
			SecsSpent: 5 * secsInOneHour,
		},
		{
			ID:        2,
			Summary:   "seeded task 2",
			Active:    true,
			CreatedAt: ca,
			UpdatedAt: ca.Add(time.Hour * 6),
			SecsSpent: 4 * secsInOneHour,
		},
	}

	taskLogs := []types.TaskLogEntry{
		{
			ID:        1,
			TaskID:    1,
			BeginTS:   ca.Add(time.Hour * 2),
			EndTS:     ca.Add(time.Hour * 4),
			SecsSpent: 2 * secsInOneHour,
			Comment:   "task 1 tl 1",
		},
		{
			ID:        2,
			TaskID:    1,
			BeginTS:   ca.Add(time.Hour * 6),
			EndTS:     ca.Add(time.Hour * 9),
			SecsSpent: 3 * secsInOneHour,
			Comment:   "task 1 tl 2",
		},
		{
			ID:        3,
			TaskID:    2,
			BeginTS:   ca.Add(time.Hour * 2),
			EndTS:     ca.Add(time.Hour * 6),
			SecsSpent: 4 * secsInOneHour,
			Comment:   "task 2 tl 1",
		},
	}

	return testData{tasks, taskLogs}
}

func seedDB(t *testing.T, db *sql.DB, data testData) {
	t.Helper()

	for _, task := range data.tasks {
		_, err := db.Exec(`
INSERT INTO task (id, summary, secs_spent, active, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, task.ID, task.Summary, task.SecsSpent, task.Active, task.CreatedAt, task.UpdatedAt)
		require.NoError(t, err, "failed to insert data into table \"task\": %v", err)
	}

	for _, taskLog := range data.taskLogs {
		_, err := db.Exec(`
INSERT INTO task_log (id, task_id, begin_ts, end_ts, secs_spent, comment, active)
VALUES (?, ?, ?, ?, ?, ?, ?)`, taskLog.ID, taskLog.TaskID, taskLog.BeginTS, taskLog.EndTS, taskLog.SecsSpent, taskLog.Comment, false)
		require.NoError(t, err, "failed to insert data into table \"task_log\": %v", err)
	}
}
