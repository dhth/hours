package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dhth/hours/internal/types"
)

func insertNewTLInDB(db *sql.DB, taskID int, beginTs time.Time) error {
	stmt, err := db.Prepare(`
INSERT INTO task_log (task_id, begin_ts, active)
VALUES (?, ?, ?);
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskID, beginTs.UTC(), true)
	if err != nil {
		return err
	}

	return nil
}

func updateTLBeginTSInDB(db *sql.DB, beginTs time.Time) error {
	stmt, err := db.Prepare(`
UPDATE task_log SET begin_ts=?
WHERE active is true;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(beginTs.UTC(), true)
	if err != nil {
		return err
	}

	return nil
}

func deleteActiveTLInDB(db *sql.DB) error {
	stmt, err := db.Prepare(`
DELETE FROM task_log
WHERE active=true;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()

	return err
}

func updateActiveTLInDB(db *sql.DB, taskLogID int, taskID int, beginTs, endTs time.Time, secsSpent int, comment string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(`
UPDATE task_log
SET active = 0,
    begin_ts = ?,
    end_ts = ?,
    secs_spent = ?,
    comment = ?
WHERE id = ?
AND active = 1;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(beginTs.UTC(), endTs.UTC(), secsSpent, comment, taskLogID)
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
UPDATE task
SET secs_spent = secs_spent+?,
    updated_at = ?
WHERE id = ?;
    `)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(secsSpent, time.Now().UTC(), taskID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func insertManualTLInDB(db *sql.DB, taskID int, beginTs time.Time, endTs time.Time, comment string) error {
	secsSpent := int(endTs.Sub(beginTs).Seconds())
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(`
INSERT INTO task_log (task_id, begin_ts, end_ts, secs_spent, comment, active)
VALUES (?, ?, ?, ?, ?, ?);
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskID, beginTs.UTC(), endTs.UTC(), secsSpent, comment, false)
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
UPDATE task
SET secs_spent = secs_spent+?,
    updated_at = ?
WHERE id = ?;
    `)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(secsSpent, time.Now().UTC(), taskID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func fetchActiveTaskFromDB(db *sql.DB) (types.ActiveTaskDetails, error) {
	row := db.QueryRow(`
SELECT t.id, t.summary, tl.begin_ts
FROM task_log tl left join task t on tl.task_id = t.id
WHERE tl.active=true;
`)

	var activeTaskDetails types.ActiveTaskDetails
	err := row.Scan(
		&activeTaskDetails.TaskID,
		&activeTaskDetails.TaskSummary,
		&activeTaskDetails.LastLogEntryBeginTS,
	)
	if errors.Is(err, sql.ErrNoRows) {
		activeTaskDetails.TaskID = -1
		return activeTaskDetails, nil
	} else if err != nil {
		return activeTaskDetails, err
	}
	activeTaskDetails.LastLogEntryBeginTS = activeTaskDetails.LastLogEntryBeginTS.Local()
	return activeTaskDetails, nil
}

func insertTaskInDB(db *sql.DB, summary string) error {
	stmt, err := db.Prepare(`
INSERT into task (summary, active, created_at, updated_at)
VALUES (?, true, ?, ?);
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC()
	_, err = stmt.Exec(summary, now, now)
	if err != nil {
		return err
	}
	return nil
}

func updateTaskInDB(db *sql.DB, id int, summary string) error {
	stmt, err := db.Prepare(`
UPDATE task
SET summary = ?,
    updated_at = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	return nil
}

func updateTaskActiveStatusInDB(db *sql.DB, id int, active bool) error {
	stmt, err := db.Prepare(`
UPDATE task
SET active = ?,
    updated_at = ?
WHERE id = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(active, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	return nil
}

func updateTaskDataFromDB(db *sql.DB, t *types.Task) error {
	row := db.QueryRow(`
SELECT secs_spent, updated_at
FROM task
WHERE id=?;
    `, t.ID)

	err := row.Scan(
		&t.SecsSpent,
		&t.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func fetchTasksFromDB(db *sql.DB, active bool, limit int) ([]types.Task, error) {
	var tasks []types.Task

	rows, err := db.Query(`
SELECT id, summary, secs_spent, created_at, updated_at, active
FROM task
WHERE active=?
ORDER by updated_at DESC
LIMIT ?;
    `, active, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry types.Task
		err = rows.Scan(&entry.ID,
			&entry.Summary,
			&entry.SecsSpent,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.Active,
		)
		if err != nil {
			return nil, err
		}
		entry.CreatedAt = entry.CreatedAt.Local()
		entry.UpdatedAt = entry.UpdatedAt.Local()
		tasks = append(tasks, entry)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return tasks, nil
}

func fetchTLEntriesFromDB(db *sql.DB, desc bool, limit int) ([]types.TaskLogEntry, error) {
	var logEntries []types.TaskLogEntry

	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}
	query := fmt.Sprintf(`
SELECT tl.id, tl.task_id, t.summary, tl.begin_ts, tl.end_ts, tl.secs_spent, tl.comment
FROM task_log tl left join task t on tl.task_id=t.id
WHERE tl.active=false
ORDER by tl.begin_ts %s
LIMIT ?;
`, order)

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry types.TaskLogEntry
		err = rows.Scan(&entry.ID,
			&entry.TaskID,
			&entry.TaskSummary,
			&entry.BeginTS,
			&entry.EndTS,
			&entry.SecsSpent,
			&entry.Comment,
		)
		if err != nil {
			return nil, err
		}
		entry.BeginTS = entry.BeginTS.Local()
		entry.EndTS = entry.EndTS.Local()
		logEntries = append(logEntries, entry)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return logEntries, nil
}

func fetchTLEntriesBetweenTSFromDB(db *sql.DB, beginTs, endTs time.Time, limit int) ([]types.TaskLogEntry, error) {
	var logEntries []types.TaskLogEntry

	rows, err := db.Query(`
SELECT tl.id, tl.task_id, t.summary, tl.begin_ts, tl.end_ts, tl.secs_spent, tl.comment
FROM task_log tl left join task t on tl.task_id=t.id
WHERE tl.active=false
AND tl.end_ts >= ?
AND tl.end_ts < ?
ORDER by tl.begin_ts ASC LIMIT ?;
    `, beginTs.UTC(), endTs.UTC(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry types.TaskLogEntry
		err = rows.Scan(&entry.ID,
			&entry.TaskID,
			&entry.TaskSummary,
			&entry.BeginTS,
			&entry.EndTS,
			&entry.SecsSpent,
			&entry.Comment,
		)
		if err != nil {
			return nil, err
		}
		entry.BeginTS = entry.BeginTS.Local()
		entry.EndTS = entry.EndTS.Local()
		logEntries = append(logEntries, entry)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return logEntries, nil
}

func fetchStatsFromDB(db *sql.DB, limit int) ([]types.TaskReportEntry, error) {
	rows, err := db.Query(`
SELECT tl.task_id, t.summary, COUNT(tl.id) as num_entries, t.secs_spent
from task_log tl
LEFT JOIN task t on tl.task_id = t.id
GROUP BY tl.task_id
ORDER BY t.secs_spent DESC
limit ?;
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tLE []types.TaskReportEntry

	for rows.Next() {
		var entry types.TaskReportEntry
		err = rows.Scan(
			&entry.TaskID,
			&entry.TaskSummary,
			&entry.NumEntries,
			&entry.SecsSpent,
		)
		if err != nil {
			return nil, err
		}
		tLE = append(tLE, entry)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return tLE, nil
}

func fetchStatsBetweenTSFromDB(db *sql.DB, beginTs, endTs time.Time, limit int) ([]types.TaskReportEntry, error) {
	rows, err := db.Query(`
SELECT tl.task_id, t.summary, COUNT(tl.id) as num_entries,  SUM(tl.secs_spent) AS secs_spent
FROM task_log tl 
LEFT JOIN task t ON tl.task_id = t.id
WHERE tl.end_ts >= ? AND tl.end_ts < ?
GROUP BY tl.task_id
ORDER BY secs_spent DESC
LIMIT ?;
`, beginTs.UTC(), endTs.UTC(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tLE []types.TaskReportEntry

	for rows.Next() {
		var entry types.TaskReportEntry
		err = rows.Scan(
			&entry.TaskID,
			&entry.TaskSummary,
			&entry.NumEntries,
			&entry.SecsSpent,
		)
		if err != nil {
			return nil, err
		}
		tLE = append(tLE, entry)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return tLE, nil
}

func fetchReportBetweenTSFromDB(db *sql.DB, beginTs, endTs time.Time, limit int) ([]types.TaskReportEntry, error) {
	rows, err := db.Query(`
SELECT tl.task_id, t.summary, COUNT(tl.id) as num_entries,  SUM(tl.secs_spent) AS secs_spent
FROM task_log tl 
LEFT JOIN task t ON tl.task_id = t.id
WHERE tl.end_ts >= ? AND tl.end_ts < ?
GROUP BY tl.task_id
ORDER BY t.updated_at ASC
LIMIT ?;
`, beginTs.UTC(), endTs.UTC(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tLE []types.TaskReportEntry

	for rows.Next() {
		var entry types.TaskReportEntry
		err = rows.Scan(
			&entry.TaskID,
			&entry.TaskSummary,
			&entry.NumEntries,
			&entry.SecsSpent,
		)
		if err != nil {
			return nil, err
		}
		tLE = append(tLE, entry)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return tLE, nil
}

func deleteEntry(db *sql.DB, entry *types.TaskLogEntry) error {
	secsSpent := entry.SecsSpent

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(`
DELETE from task_log
WHERE ID=?;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entry.ID)
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
UPDATE task
SET secs_spent = secs_spent-?,
    updated_at = ?
WHERE id = ?;
    `)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(secsSpent, time.Now().UTC(), entry.TaskID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
