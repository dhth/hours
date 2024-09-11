package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dhth/hours/internal/types"
)

var ErrCouldntRollBackTx = errors.New("couldn't roll back transaction")

func InsertNewTL(db *sql.DB, taskID int, beginTs time.Time) (int, error) {
	return runInTxAndReturnID(db, func(tx *sql.Tx) (int, error) {
		stmt, err := tx.Prepare(`
INSERT INTO task_log (task_id, begin_ts, active)
VALUES (?, ?, ?);
`)
		if err != nil {
			return -1, err
		}
		defer stmt.Close()

		res, err := stmt.Exec(taskID, beginTs.UTC(), true)
		if err != nil {
			return -1, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return -1, err
		}

		return int(lastID), nil
	})
}

func UpdateTLBeginTS(db *sql.DB, beginTs time.Time) error {
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

func DeleteActiveTL(db *sql.DB) error {
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

func UpdateActiveTL(db *sql.DB, taskLogID int, taskID int, beginTs, endTs time.Time, secsSpent int, comment string) error {
	return runInTx(db, func(tx *sql.Tx) error {
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

		return err
	})
}

func InsertManualTL(db *sql.DB, taskID int, beginTs time.Time, endTs time.Time, comment string) (int, error) {
	return runInTxAndReturnID(db, func(tx *sql.Tx) (int, error) {
		stmt, err := tx.Prepare(`
INSERT INTO task_log (task_id, begin_ts, end_ts, secs_spent, comment, active)
VALUES (?, ?, ?, ?, ?, ?);
`)
		if err != nil {
			return -1, err
		}
		defer stmt.Close()

		secsSpent := int(endTs.Sub(beginTs).Seconds())

		res, err := stmt.Exec(taskID, beginTs.UTC(), endTs.UTC(), secsSpent, comment, false)
		if err != nil {
			return -1, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return -1, err
		}

		tStmt, err := tx.Prepare(`
UPDATE task
SET secs_spent = secs_spent+?,
    updated_at = ?
WHERE id = ?;
    `)
		if err != nil {
			return -1, err
		}
		defer tStmt.Close()

		_, err = tStmt.Exec(secsSpent, time.Now().UTC(), taskID)
		if err != nil {
			return -1, err
		}

		return int(lastID), nil
	})
}

func FetchActiveTask(db *sql.DB) (types.ActiveTaskDetails, error) {
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

func InsertTask(db *sql.DB, summary string) (int, error) {
	return runInTxAndReturnID(db, func(tx *sql.Tx) (int, error) {
		stmt, err := tx.Prepare(`
INSERT into task (summary, active, created_at, updated_at)
VALUES (?, true, ?, ?);
`)
		if err != nil {
			return -1, err
		}
		defer stmt.Close()

		now := time.Now().UTC()
		res, err := stmt.Exec(summary, now, now)
		if err != nil {
			return -1, err
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return -1, err
		}

		return int(lastID), nil
	})
}

func UpdateTask(db *sql.DB, id int, summary string) error {
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

func UpdateTaskActiveStatus(db *sql.DB, id int, active bool) error {
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

func UpdateTaskData(db *sql.DB, t *types.Task) error {
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

func FetchTasks(db *sql.DB, active bool, limit int) ([]types.Task, error) {
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

func FetchTLEntries(db *sql.DB, desc bool, limit int) ([]types.TaskLogEntry, error) {
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

func FetchTLEntriesBetweenTS(db *sql.DB, beginTs, endTs time.Time, limit int) ([]types.TaskLogEntry, error) {
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

func FetchStats(db *sql.DB, limit int) ([]types.TaskReportEntry, error) {
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

func FetchStatsBetweenTS(db *sql.DB, beginTs, endTs time.Time, limit int) ([]types.TaskReportEntry, error) {
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

func FetchReportBetweenTS(db *sql.DB, beginTs, endTs time.Time, limit int) ([]types.TaskReportEntry, error) {
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

func DeleteTaskLogEntry(db *sql.DB, entry *types.TaskLogEntry) error {
	return runInTx(db, func(tx *sql.Tx) error {
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

		_, err = tStmt.Exec(entry.SecsSpent, time.Now().UTC(), entry.TaskID)
		return err
	})
}

func runInTxAndReturnID(db *sql.DB, fn func(tx *sql.Tx) (int, error)) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	lastID, err := fn(tx)
	if err == nil {
		return lastID, tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return lastID, fmt.Errorf("%w: %w: %s", ErrCouldntRollBackTx, rollbackErr, err.Error())
	}

	return lastID, err
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return fmt.Errorf("%w: %w: %w", ErrCouldntRollBackTx, rollbackErr, err)
	}

	return err
}

func fetchTaskByID(db *sql.DB, id int) (types.Task, error) {
	var task types.Task
	row := db.QueryRow(`
SELECT id, summary, secs_spent, active, created_at, updated_at
FROM task
WHERE id=?;
    `, id)

	if row.Err() != nil {
		return task, row.Err()
	}
	err := row.Scan(&task.ID,
		&task.Summary,
		&task.SecsSpent,
		&task.Active,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return task, err
	}
	task.CreatedAt = task.CreatedAt.Local()
	task.UpdatedAt = task.UpdatedAt.Local()

	return task, nil
}

func fetchTaskLogByID(db *sql.DB, id int) (types.TaskLogEntry, error) {
	var tl types.TaskLogEntry
	row := db.QueryRow(`
SELECT id, task_id, begin_ts, end_ts, secs_spent, comment
FROM task_log
WHERE id=?;
    `, id)

	if row.Err() != nil {
		return tl, row.Err()
	}
	err := row.Scan(&tl.ID,
		&tl.TaskID,
		&tl.BeginTS,
		&tl.EndTS,
		&tl.SecsSpent,
		&tl.Comment,
	)
	if err != nil {
		return tl, err
	}
	tl.BeginTS = tl.BeginTS.Local()
	tl.EndTS = tl.EndTS.Local()

	return tl, nil
}
