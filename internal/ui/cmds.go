package ui

import (
	"database/sql"
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/types"
	_ "modernc.org/sqlite" // sqlite driver
)

func toggleTracking(db *sql.DB,
	taskID int,
	beginTs time.Time,
	endTs time.Time,
	comment string,
) tea.Cmd {
	return func() tea.Msg {
		row := db.QueryRow(`
SELECT id, task_id
FROM task_log
WHERE active=1
ORDER BY begin_ts DESC
LIMIT 1
`)
		var trackStatus trackingStatus
		var activeTaskLogID int
		var activeTaskID int

		err := row.Scan(&activeTaskLogID, &activeTaskID)
		if errors.Is(err, sql.ErrNoRows) {
			trackStatus = trackingInactive
		} else if err != nil {
			return trackingToggledMsg{err: err}
		} else {
			trackStatus = trackingActive
		}

		switch trackStatus {
		case trackingInactive:
			_, err = pers.InsertNewTL(db, taskID, beginTs)
			if err != nil {
				return trackingToggledMsg{err: err}
			} else {
				return trackingToggledMsg{taskID: taskID}
			}

		default:
			secsSpent := int(endTs.Sub(beginTs).Seconds())
			err := pers.UpdateActiveTL(db, activeTaskLogID, activeTaskID, beginTs, endTs, secsSpent, comment)
			if err != nil {
				return trackingToggledMsg{err: err}
			} else {
				return trackingToggledMsg{taskID: taskID, finished: true, secsSpent: secsSpent}
			}
		}
	}
}

func updateTLBeginTS(db *sql.DB, beginTS time.Time) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTLBeginTS(db, beginTS)
		return tlBeginTSUpdatedMsg{beginTS, err}
	}
}

func insertManualTL(db *sql.DB, taskID int, beginTS time.Time, endTS time.Time, comment string) tea.Cmd {
	return func() tea.Msg {
		_, err := pers.InsertManualTL(db, taskID, beginTS, endTS, comment)
		return manualTLInsertedMsg{taskID, err}
	}
}

func fetchActiveTask(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		activeTaskDetails, err := pers.FetchActiveTask(db)
		if err != nil {
			return activeTaskFetchedMsg{err: err}
		}

		if activeTaskDetails.TaskID == -1 {
			return activeTaskFetchedMsg{noneActive: true}
		}

		return activeTaskFetchedMsg{
			activeTaskID: activeTaskDetails.TaskID,
			beginTs:      activeTaskDetails.LastLogEntryBeginTS,
		}
	}
}

func updateTaskRep(db *sql.DB, t *types.Task) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTaskData(db, t)
		return taskRepUpdatedMsg{
			tsk: t,
			err: err,
		}
	}
}

func fetchTLS(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		entries, err := pers.FetchTLEntries(db, true, 50)
		return tLsFetchedMsg{
			entries: entries,
			err:     err,
		}
	}
}

func deleteTL(db *sql.DB, entry *types.TaskLogEntry) tea.Cmd {
	return func() tea.Msg {
		err := pers.DeleteTL(db, entry)
		return tLDeletedMsg{
			entry: entry,
			err:   err,
		}
	}
}

func deleteActiveTL(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		err := pers.DeleteActiveTL(db)
		return activeTaskLogDeletedMsg{err}
	}
}

func createTask(db *sql.DB, summary string) tea.Cmd {
	return func() tea.Msg {
		_, err := pers.InsertTask(db, summary)
		return taskCreatedMsg{err}
	}
}

func updateTask(db *sql.DB, task *types.Task, summary string) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTask(db, task.ID, summary)
		return taskUpdatedMsg{task, summary, err}
	}
}

func updateTaskActiveStatus(db *sql.DB, task *types.Task, active bool) tea.Cmd {
	return func() tea.Msg {
		err := pers.UpdateTaskActiveStatus(db, task.ID, active)
		return taskActiveStatusUpdatedMsg{task, active, err}
	}
}

func fetchTasks(db *sql.DB, active bool) tea.Cmd {
	return func() tea.Msg {
		tasks, err := pers.FetchTasks(db, active, 50)
		return tasksFetchedMsg{tasks, active, err}
	}
}

func hideHelp(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return hideHelpMsg{}
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
			data, err = getTaskLog(db, start, end, 20, plain)
		case reportStats:
			data, err = getStats(db, period, start, end, plain)
		}

		return recordsDataFetchedMsg{
			start:  start,
			end:    end,
			report: data,
			err:    err,
		}
	}
}
