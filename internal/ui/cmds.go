package ui

import (
	"database/sql"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	_ "modernc.org/sqlite"
)

func toggleTracking(db *sql.DB,
	taskId int,
	beginTs time.Time,
	endTs time.Time,
	comment string) tea.Cmd {
	return func() tea.Msg {

		row := db.QueryRow(`
SELECT id, task_id
FROM task_log
WHERE active=1
ORDER BY begin_ts DESC
LIMIT 1
`)
		var trackStatus trackingStatus
		var activeTaskLogId int
		var activeTaskId int

		err := row.Scan(&activeTaskLogId, &activeTaskId)
		if err == sql.ErrNoRows {
			trackStatus = trackingInactive
		} else if err != nil {
			return trackingToggledMsg{err: err}
		} else {
			trackStatus = trackingActive
		}

		switch trackStatus {
		case trackingInactive:
			err = insertNewTLInDB(db, taskId, beginTs)
			if err != nil {
				return trackingToggledMsg{err: err}
			} else {
				return trackingToggledMsg{taskId: taskId}
			}

		default:
			secsSpent := int(endTs.Sub(beginTs).Seconds())
			err := updateActiveTLInDB(db, activeTaskLogId, activeTaskId, endTs, secsSpent, comment)
			if err != nil {
				return trackingToggledMsg{err: err}
			} else {
				return trackingToggledMsg{taskId: taskId, finished: true, secsSpent: secsSpent}
			}
		}
	}
}

func insertManualEntry(db *sql.DB, taskId int, beginTS time.Time, endTS time.Time, comment string) tea.Cmd {
	return func() tea.Msg {
		err := insertManualTLInDB(db, taskId, beginTS, endTS, comment)
		return manualTaskLogInserted{taskId, err}
	}
}

func fetchActiveTask(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		id, beginTs, err := fetchActiveTaskFromDB(db)

		if err != nil {
			return activeTaskFetchedMsg{err: err}
		}

		if id == -1 {
			return activeTaskFetchedMsg{noneActive: true}
		}

		return activeTaskFetchedMsg{activeTaskId: id, beginTs: beginTs}
	}
}

func updateTaskRep(db *sql.DB, t *task) tea.Cmd {
	return func() tea.Msg {
		err := updateTaskDataFromDB(db, t)
		return taskRepUpdatedMsg{
			tsk: t,
			err: err,
		}
	}
}

func fetchTaskLogEntries(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		entries, err := fetchTLEntriesFromDB(db)
		return taskLogEntriesFetchedMsg{
			entries: entries,
			err:     err,
		}
	}
}

func deleteLogEntry(db *sql.DB, entry *taskLogEntry) tea.Cmd {
	return func() tea.Msg {
		err := deleteEntry(db, entry)
		return taskLogEntryDeletedMsg{
			entry: entry,
			err:   err,
		}
	}
}

func createTask(db *sql.DB, summary string) tea.Cmd {
	return func() tea.Msg {
		err := insertTaskInDB(db, summary)
		return taskCreatedMsg{err}
	}
}

func updateTask(db *sql.DB, task *task, summary string) tea.Cmd {
	return func() tea.Msg {
		err := updateTaskInDB(db, task.id, summary)
		return taskUpdatedMsg{task, summary, err}
	}
}

func fetchTasks(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		tasks, err := fetchTasksFromDB(db)
		return tasksFetched{tasks, err}
	}
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return HideHelpMsg{}
	})
}
