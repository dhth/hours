package ui

import (
	"database/sql"
	"errors"
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

	_, err = stmt.Exec(taskId, beginTs, true)
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
    comment = ?
WHERE id = ?
AND active = 1;
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(beginTs, endTs, comment, taskLogId)
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
UPDATE task
SET secsSpent = secsSpent+?,
    updated_at = ?
WHERE id = ?;
    `)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(secsSpent, time.Now().Local(), taskId)
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
INSERT INTO task_log (task_id, begin_ts, end_ts, comment, active)
VALUES (?, ?, ?, ?, ?);
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskId, beginTs, endTs, comment, false)
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
UPDATE task
SET secsSpent = secsSpent+?,
    updated_at = ?
WHERE id = ?;
    `)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(secsSpent, time.Now().Local(), taskId)
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
		&activeTaskDetails.lastLogEntryBeginTS,
	)
	if errors.Is(err, sql.ErrNoRows) {
		activeTaskDetails.taskId = -1
		return activeTaskDetails, nil
	} else if err != nil {
		return activeTaskDetails, err
	}
	return activeTaskDetails, nil
}

func insertTaskInDB(db *sql.DB, summary string) error {

	stmt, err := db.Prepare(`
INSERT into task (summary, active, created_at, updated_at)
VALUES (?, true, DATETIME('now'), DATETIME('now'));
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(summary)

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

	_, err = stmt.Exec(summary, time.Now().Local(), id)

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

	_, err = stmt.Exec(active, time.Now().Local(), id)

	if err != nil {
		return err
	}
	return nil
}

func updateTaskDataFromDB(db *sql.DB, t *task) error {

	row := db.QueryRow(`
SELECT secsSpent, updated_at
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
SELECT id, summary, secsSpent, created_at, updated_at, active
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
		tasks = append(tasks, entry)

	}
	return tasks, nil
}

func fetchTLEntriesFromDB(db *sql.DB, limit int) ([]taskLogEntry, error) {

	var logEntries []taskLogEntry

	rows, err := db.Query(`
SELECT tl.id, tl.task_id, t.summary, tl.begin_ts, tl.end_ts, tl.comment
FROM task_log tl left join task t on tl.task_id=t.id
WHERE tl.active=false
ORDER by tl.begin_ts DESC LIMIT ?;
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry taskLogEntry
		err = rows.Scan(&entry.id,
			&entry.taskId,
			&entry.taskSummary,
			&entry.beginTS,
			&entry.endTS,
			&entry.comment,
		)
		if err != nil {
			return nil, err
		}
		logEntries = append(logEntries, entry)

	}
	return logEntries, nil
}

func fetchTLEntriesFromDBAfterTS(db *sql.DB, beginTs time.Time, limit int) ([]taskLogEntry, error) {

	var logEntries []taskLogEntry

	rows, err := db.Query(`
SELECT tl.id, tl.task_id, t.summary, tl.begin_ts, tl.end_ts, tl.comment
FROM task_log tl left join task t on tl.task_id=t.id
WHERE tl.active=false
AND begin_ts >= ?
ORDER by tl.begin_ts ASC LIMIT ?;
    `, beginTs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry taskLogEntry
		err = rows.Scan(&entry.id,
			&entry.taskId,
			&entry.taskSummary,
			&entry.beginTS,
			&entry.endTS,
			&entry.comment,
		)
		if err != nil {
			return nil, err
		}
		logEntries = append(logEntries, entry)

	}
	return logEntries, nil
}

func deleteEntry(db *sql.DB, entry *taskLogEntry) error {
	secsSpent := int(entry.endTS.Sub(entry.beginTS).Seconds())

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
SET secsSpent = secsSpent-?,
    updated_at = ?
WHERE id = ?;
    `)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(secsSpent, time.Now().Local(), entry.taskId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
