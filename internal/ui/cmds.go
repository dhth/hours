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
	comment *string,
) tea.Cmd {
	return func() tea.Msg {
		row := db.QueryRow(`
SELECT id, task_id
FROM task_log
WHERE active=1
ORDER BY begin_ts DESC
LIMIT 1
`)
		var isTrackingActive bool
		var activeTaskLogID int
		var activeTaskID int

		err := row.Scan(&activeTaskLogID, &activeTaskID)
		if errors.Is(err, sql.ErrNoRows) {
			isTrackingActive = false
		} else if err != nil {
			return trackingToggledMsg{err: err}
		} else {
			isTrackingActive = true
		}

		switch isTrackingActive {
		case false:
			_, err = pers.InsertNewTL(db, taskID, beginTs)
			if err != nil {
				return trackingToggledMsg{err: err}
			}
			return trackingToggledMsg{taskID: taskID}

		default:
			secsSpent := int(endTs.Sub(beginTs).Seconds())
			err := pers.FinishActiveTL(db, activeTaskLogID, activeTaskID, beginTs, endTs, secsSpent, comment)
			if err != nil {
				return trackingToggledMsg{err: err}
			}
			return trackingToggledMsg{taskID: taskID, finished: true, secsSpent: secsSpent}
		}
	}
}

func quickSwitchActiveIssue(db *sql.DB, taskID int, ts time.Time) tea.Cmd {
	return func() tea.Msg {
		result, err := pers.QuickSwitchActiveTL(db, taskID, ts)
		return activeTLSwitchedMsg{
			lastActiveTaskID:      result.LastActiveTaskID,
			currentlyActiveTaskID: taskID,
			currentlyActiveTLID:   result.CurrentlyActiveTLID,
			ts:                    ts,
			err:                   err,
		}
	}
}

func updateActiveTL(db *sql.DB, beginTS time.Time, comment *string) tea.Cmd {
	return func() tea.Msg {
		err := pers.EditActiveTL(db, beginTS, comment)
		return activeTLUpdatedMsg{beginTS, comment, err}
	}
}

func insertManualTL(db *sql.DB, taskID int, beginTS time.Time, endTS time.Time, comment *string) tea.Cmd {
	return func() tea.Msg {
		_, err := pers.InsertManualTL(db, taskID, beginTS, endTS, comment)
		return manualTLInsertedMsg{taskID, err}
	}
}

func editSavedTL(db *sql.DB, tlID, taskID int, beginTS time.Time, endTS time.Time, comment *string) tea.Cmd {
	return func() tea.Msg {
		_, err := pers.EditSavedTL(db, tlID, beginTS, endTS, comment)
		return savedTLEditedMsg{tlID, taskID, err}
	}
}

func fetchActiveTask(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		activeTaskDetails, err := pers.FetchActiveTaskDetails(db)
		if err != nil {
			return activeTaskFetchedMsg{err: err}
		}

		if activeTaskDetails.TaskID == -1 {
			return activeTaskFetchedMsg{noneActive: true}
		}

		return activeTaskFetchedMsg{
			activeTask: activeTaskDetails,
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

func fetchTLS(db *sql.DB, tlIDToFocusOn *int) tea.Cmd {
	return func() tea.Msg {
		entries, err := pers.FetchTLEntries(db, true, 50)
		return tLsFetchedMsg{
			entries:       entries,
			tlIDToFocusOn: tlIDToFocusOn,
			err:           err,
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

func getRecordsData(
	analyticsType recordsKind,
	db *sql.DB,
	style Style,
	dateRange types.DateRange,
	taskStatus types.TaskStatus,
	plain bool,
) tea.Cmd {
	return func() tea.Msg {
		var data string
		var err error

		switch analyticsType {
		case reportRecords:
			data, err = getReport(db, style, dateRange.Start, dateRange.NumDays, taskStatus, plain)
		case reportAggRecords:
			data, err = getReportAgg(db, style, dateRange.Start, dateRange.NumDays, taskStatus, plain)
		case reportLogs:
			data, err = getTaskLog(db, style, dateRange.Start, dateRange.End, taskStatus, 20, plain)
		case reportStats:
			data, err = getStats(db, style, &dateRange, taskStatus, plain)
		}

		return recordsDataFetchedMsg{
			dateRange: dateRange,
			report:    data,
			err:       err,
		}
	}
}

func moveTaskLog(db *sql.DB, tlID int, oldTaskID int, newTaskID int, secsSpent int) tea.Cmd {
	return func() tea.Msg {
		err := pers.MoveTaskLog(db, tlID, oldTaskID, newTaskID, secsSpent)
		return taskLogMovedMsg{tlID, oldTaskID, newTaskID, err}
	}
}
