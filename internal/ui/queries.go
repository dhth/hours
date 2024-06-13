package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func insertNewTLInDB(db *sql.DB, taskId int, beginTs time.Time) error {

	stmt, err := db.Prepare(`
    INSERT INTO task_log (task_id, begin_ts, active)
    VALUES (?, ?, ?);
    `)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskId, beginTs.UTC(), true)
	if err != nil {
		return err
	}

	return nil
}

func updateActiveTLInDB(db *sql.DB, taskLogId int, taskId int, beginTs, endTs time.Time, secsSpent int, comment string) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

	_, err = stmt.Exec(beginTs.UTC(), endTs.UTC(), secsSpent, comment, taskLogId)
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

	_, err = tStmt.Exec(secsSpent, time.Now().UTC(), taskId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func insertManualTLInDB(db *sql.DB, taskId int, beginTs time.Time, endTs time.Time, comment string) error {

	secsSpent := int(endTs.Sub(beginTs).Seconds())
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
INSERT INTO task_log (task_id, begin_ts, end_ts, secs_spent, comment, active)
VALUES (?, ?, ?, ?, ?, ?);
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskId, beginTs.UTC(), endTs.UTC(), secsSpent, comment, false)
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

	_, err = tStmt.Exec(secsSpent, time.Now().UTC(), taskId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func fetchActiveTaskFromDB(db *sql.DB) (activeTaskDetails, error) {

	row := db.QueryRow(`
SELECT t.id, t.summary, tl.begin_ts
FROM task_log tl left join task t on tl.task_id = t.id
WHERE tl.active=true;
`)

	var activeTaskDetails activeTaskDetails
	err := row.Scan(
		&activeTaskDetails.taskId,
		&activeTaskDetails.taskSummary,
		&activeTaskDetails.lastLogEntryBeginTs,
	)
	if errors.Is(err, sql.ErrNoRows) {
		activeTaskDetails.taskId = -1
		return activeTaskDetails, nil
	} else if err != nil {
		return activeTaskDetails, err
	}
	activeTaskDetails.lastLogEntryBeginTs = activeTaskDetails.lastLogEntryBeginTs.Local()
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

func updateTaskDataFromDB(db *sql.DB, t *task) error {

	row := db.QueryRow(`
SELECT secs_spent, updated_at
FROM task
WHERE id=?;
    `, t.id)

	err := row.Scan(
		&t.secsSpent,
		&t.updatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func fetchTasksFromDB(db *sql.DB, active bool, limit int) ([]task, error) {

	var tasks []task

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
		var entry task
		err = rows.Scan(&entry.id,
			&entry.summary,
			&entry.secsSpent,
			&entry.createdAt,
			&entry.updatedAt,
			&entry.active,
		)
		if err != nil {
			return nil, err
		}
		entry.createdAt = entry.createdAt.Local()
		entry.updatedAt = entry.updatedAt.Local()
		tasks = append(tasks, entry)

	}
	return tasks, nil
}

func fetchTLEntriesFromDB(db *sql.DB, desc bool, limit int) ([]taskLogEntry, error) {

	var logEntries []taskLogEntry

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
		var entry taskLogEntry
		err = rows.Scan(&entry.id,
			&entry.taskId,
			&entry.taskSummary,
			&entry.beginTs,
			&entry.endTs,
			&entry.secsSpent,
			&entry.comment,
		)
		if err != nil {
			return nil, err
		}
		entry.beginTs = entry.beginTs.Local()
		entry.endTs = entry.endTs.Local()
		logEntries = append(logEntries, entry)

	}
	return logEntries, nil
}

func fetchTLEntriesBetweenTSFromDB(db *sql.DB, beginTs, endTs time.Time, limit int) ([]taskLogEntry, error) {

	var logEntries []taskLogEntry

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
		var entry taskLogEntry
		err = rows.Scan(&entry.id,
			&entry.taskId,
			&entry.taskSummary,
			&entry.beginTs,
			&entry.endTs,
			&entry.secsSpent,
			&entry.comment,
		)
		if err != nil {
			return nil, err
		}
		entry.beginTs = entry.beginTs.Local()
		entry.endTs = entry.endTs.Local()
		logEntries = append(logEntries, entry)

	}
	return logEntries, nil
}

func fetchStatsFromDB(db *sql.DB, limit int) ([]taskReportEntry, error) {

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

	var tLE []taskReportEntry

	for rows.Next() {
		var entry taskReportEntry
		err = rows.Scan(
			&entry.taskId,
			&entry.taskSummary,
			&entry.numEntries,
			&entry.secsSpent,
		)
		if err != nil {
			return nil, err
		}
		tLE = append(tLE, entry)

	}
	return tLE, nil
}

func fetchStatsBetweenTSFromDB(db *sql.DB, beginTs, endTs time.Time, limit int) ([]taskReportEntry, error) {

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

	var tLE []taskReportEntry

	for rows.Next() {
		var entry taskReportEntry
		err = rows.Scan(
			&entry.taskId,
			&entry.taskSummary,
			&entry.numEntries,
			&entry.secsSpent,
		)
		if err != nil {
			return nil, err
		}
		tLE = append(tLE, entry)

	}
	return tLE, nil
}

func fetchReportBetweenTSFromDB(db *sql.DB, beginTs, endTs time.Time, limit int) ([]taskReportEntry, error) {

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

	var tLE []taskReportEntry

	for rows.Next() {
		var entry taskReportEntry
		err = rows.Scan(
			&entry.taskId,
			&entry.taskSummary,
			&entry.numEntries,
			&entry.secsSpent,
		)
		if err != nil {
			return nil, err
		}
		tLE = append(tLE, entry)

	}
	return tLE, nil
}

func deleteEntry(db *sql.DB, entry *taskLogEntry) error {
	secsSpent := entry.secsSpent

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
DELETE from task_log
WHERE ID=?;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(entry.id)
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

	_, err = tStmt.Exec(secsSpent, time.Now().UTC(), entry.taskId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
