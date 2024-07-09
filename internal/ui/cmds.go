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
			err := updateActiveTLInDB(db, activeTaskLogId, activeTaskId, beginTs, endTs, secsSpent, comment)
			if err != nil {
				return trackingToggledMsg{err: err}
			} else {
				return trackingToggledMsg{taskId: taskId, finished: true, secsSpent: secsSpent}
			}
		}
	}
}

func updateTLBeginTS(db *sql.DB, beginTS time.Time) tea.Cmd {
	return func() tea.Msg {
		err := updateTLBeginTSInDB(db, beginTS)
		return tlBeginTSUpdatedMsg{beginTS, err}
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
		activeTaskDetails, err := fetchActiveTaskFromDB(db)

		if err != nil {
			return activeTaskFetchedMsg{err: err}
		}

		if activeTaskDetails.taskId == -1 {
			return activeTaskFetchedMsg{noneActive: true}
		}

		return activeTaskFetchedMsg{
			activeTaskId: activeTaskDetails.taskId,
			beginTs:      activeTaskDetails.lastLogEntryBeginTs,
		}
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
		entries, err := fetchTLEntriesFromDB(db, true, 50)
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

func deleteActiveTaskLog(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		err := deleteActiveTLInDB(db)
		return activeTaskLogDeletedMsg{err}
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

func updateTaskActiveStatus(db *sql.DB, task *task, active bool) tea.Cmd {
	return func() tea.Msg {
		err := updateTaskActiveStatusInDB(db, task.id, active)
		return taskActiveStatusUpdated{task, active, err}
	}
}

func fetchTasks(db *sql.DB, active bool) tea.Cmd {
	return func() tea.Msg {
		tasks, err := fetchTasksFromDB(db, active, 50)
		return tasksFetched{tasks, active, err}
	}
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return HideHelpMsg{}
	})
}

func getRecordsData(analyticsType recordsType, db *sql.DB, period string, start, end time.Time, numDays int, plain bool) tea.Cmd {
	return func() tea.Msg {
		var data string
		var err error

		switch analyticsType {
		case reportRecords:
			data, err = getReport(db, start, numDays, plain)
		case reportAggRecords:
			data, err = getReportAgg(db, start, numDays, plain)
		case reportLogs:
			data, err = renderTaskLog(db, start, end, 20, plain)
		case reportStats:
			data, err = renderStats(db, period, start, end, plain)
		}

		return recordsDataFetchedMsg{
			start:  start,
			end:    end,
			report: data,
			err:    err,
		}
	}
}
