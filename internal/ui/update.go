package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/hours/internal/types"
)

const (
	viewPortMoveLineCount = 3
	msgCouldntSelectATask = "Couldn't select a task"
	msgChangesLocked      = "Changes locked momentarily"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.message = ""

	keyMsg, keyMsgOK := msg.(tea.KeyMsg)
	if keyMsgOK {
		if m.activeTasksList.FilterState() == list.Filtering {
			m.activeTasksList, cmd = m.activeTasksList.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch keyMsg.String() {
		case "enter":
			var updateCmd tea.Cmd
			switch m.activeView {
			case taskInputView:
				updateCmd = m.getCmdToCreateOrUpdateTask()
			case editActiveTLView:
				updateCmd = m.getCmdToUpdateActiveTL()
			case saveActiveTLView:
				updateCmd = m.getCmdToSaveActiveTL()
			case manualTasklogEntryView:
				updateCmd = m.getCmdToSaveOrUpdateTL()
			}
			if updateCmd != nil {
				cmds = append(cmds, updateCmd)
				return m, tea.Batch(cmds...)
			}

		case "esc":
			m.handleEscape()
		case "tab":
			m.goForwardInView()
		case "shift+tab":
			m.goBackwardInView()
		case "k":
			err := m.shiftTime(types.ShiftBackward, types.ShiftMinute)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case "j":
			err := m.shiftTime(types.ShiftForward, types.ShiftMinute)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case "K":
			err := m.shiftTime(types.ShiftBackward, types.ShiftFiveMinutes)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case "J":
			err := m.shiftTime(types.ShiftForward, types.ShiftFiveMinutes)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case "h":
			err := m.shiftTime(types.ShiftBackward, types.ShiftDay)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		case "l":
			err := m.shiftTime(types.ShiftForward, types.ShiftDay)
			if err != nil {
				return m, tea.Batch(cmds...)
			}
		}
	}

	switch m.activeView {
	case taskInputView:
		for i := range m.taskInputs {
			m.taskInputs[i], cmd = m.taskInputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case editActiveTLView:
		m.trackingInputs[entryBeginTS], cmd = m.trackingInputs[entryBeginTS].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case saveActiveTLView:
		for i := range m.trackingInputs {
			m.trackingInputs[i], cmd = m.trackingInputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case manualTasklogEntryView:
		for i := range m.trackingInputs {
			m.trackingInputs[i], cmd = m.trackingInputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			shouldQuit := m.handleRequestToGoBackOrQuit()
			if shouldQuit {
				return m, tea.Quit
			}
		case "1":
			if m.activeView != taskListView {
				m.activeView = taskListView
			}
		case "2":
			if m.activeView != taskLogView {
				m.activeView = taskLogView
			}
		case "3":
			if m.activeView != inactiveTaskListView {
				m.activeView = inactiveTaskListView
			}
		case "ctrl+r":
			reloadCmd := m.getCmdToReloadData()
			if reloadCmd != nil {
				cmds = append(cmds, reloadCmd)
			}
		case "ctrl+t":
			m.goToActiveTask()
		case "ctrl+s":
			if m.activeView == taskListView {
				switch m.trackingActive {
				case true:
					m.handleRequestToSaveActiveTL()
				case false:
					m.handleRequestToCreateManualTL()
				}
			}
		case "ctrl+d":
			var handleCmd tea.Cmd
			switch m.activeView {
			case taskListView:
				handleCmd = m.getCmdToDeactivateTask()
			case taskLogView:
				handleCmd = m.getCmdToDeleteTL()
			case inactiveTaskListView:
				handleCmd = m.getCmdToActivateDeactivatedTask()
			}
			if handleCmd != nil {
				cmds = append(cmds, handleCmd)
			}
		case "ctrl+x":
			if m.activeView == taskListView && m.trackingActive {
				cmds = append(cmds, deleteActiveTL(m.db))
			}
		case "s":
			if m.activeView == taskListView {
				switch m.lastTrackingChange {
				case trackingFinished:
					trackCmd := m.getCmdToStartTracking()
					if trackCmd != nil {
						cmds = append(cmds, trackCmd)
					}
				case trackingStarted:
					m.handleRequestToStopTracking()
				}
			}
		case "a":
			if m.activeView == taskListView {
				m.handleRequestToCreateTask()
			}
		case "u":
			if m.activeView == taskListView {
				m.handleRequestToUpdateTask()
			}
		case "k":
			if m.activeView == helpView {
				m.handleRequestToScrollVPUp()
			}
		case "j":
			if m.activeView == helpView {
				m.handleRequestToScrollVPDown()
			}
		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
		}

	case tea.WindowSizeMsg:
		m.handleWindowResizing(msg)
	case taskCreatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error creating task: %s", msg.err)
		} else {
			cmds = append(cmds, fetchTasks(m.db, true))
		}
	case taskUpdatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error updating task: %s", msg.err)
		} else {
			msg.tsk.Summary = msg.summary
			msg.tsk.UpdateTitle()
		}
	case tasksFetchedMsg:
		handleCmd := m.handleTasksFetchedMsg(msg)
		if handleCmd != nil {
			cmds = append(cmds, handleCmd)
		}
	case tlBeginTSUpdatedMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			m.activeTLBeginTS = msg.beginTS
		}
	case manualTLInsertedMsg:
		handleCmds := m.handleManualTLInsertedMsg(msg)
		if handleCmds != nil {
			cmds = append(cmds, handleCmds...)
		}
	case tLsFetchedMsg:
		m.handleTLSFetchedMsg(msg)
	case activeTaskFetchedMsg:
		m.handleActiveTaskFetchedMsg(msg)
	case trackingToggledMsg:
		updateCmds := m.handleTrackingToggledMsg(msg)
		if updateCmds != nil {
			cmds = append(cmds, updateCmds...)
		}
	case taskRepUpdatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error updating task status: %s", msg.err)
		} else {
			msg.tsk.UpdateDesc()
		}
	case tLDeletedMsg:
		updateCmds := m.handleTLDeleted(msg)
		if updateCmds != nil {
			cmds = append(cmds, updateCmds...)
		}
	case activeTaskLogDeletedMsg:
		m.handleActiveTLDeletedMsg(msg)
	case taskActiveStatusUpdatedMsg:
		if msg.err != nil {
			m.message = "error updating task's active status: " + msg.err.Error()
		} else {
			cmds = append(cmds, fetchTasks(m.db, true))
			cmds = append(cmds, fetchTasks(m.db, false))
		}
	case hideHelpMsg:
		m.showHelpIndicator = false
	}

	switch m.activeView {
	case taskListView:
		m.activeTasksList, cmd = m.activeTasksList.Update(msg)
		cmds = append(cmds, cmd)
	case taskLogView:
		m.taskLogList, cmd = m.taskLogList.Update(msg)
		cmds = append(cmds, cmd)
	case inactiveTaskListView:
		m.inactiveTasksList, cmd = m.inactiveTasksList.Update(msg)
		cmds = append(cmds, cmd)
	case helpView:
		m.helpVP, cmd = m.helpVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m recordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			if !m.busy {
				var newStart, newEnd time.Time
				var numDays int

				switch m.period {
				case types.TimePeriodWeek:
					weekday := m.start.Weekday()
					offset := (7 + weekday - time.Monday) % 7
					startOfPrevWeek := m.start.AddDate(0, 0, -int(offset+7))
					newStart = time.Date(startOfPrevWeek.Year(), startOfPrevWeek.Month(), startOfPrevWeek.Day(), 0, 0, 0, 0, startOfPrevWeek.Location())
					numDays = 7
				default:
					newStart = m.start.AddDate(0, 0, -m.numDays)
					numDays = m.numDays
				}
				newEnd = newStart.AddDate(0, 0, numDays)
				cmds = append(cmds, getRecordsData(m.typ, m.db, m.period, newStart, newEnd, numDays, m.plain))
				m.busy = true
			}
		case "right", "l":
			if !m.busy {
				var newStart, newEnd time.Time
				var numDays int

				switch m.period {
				case types.TimePeriodWeek:
					weekday := m.start.Weekday()
					offset := (7 + weekday - time.Monday) % 7
					startOfNextWeek := m.start.AddDate(0, 0, 7-int(offset))
					newStart = time.Date(startOfNextWeek.Year(), startOfNextWeek.Month(), startOfNextWeek.Day(), 0, 0, 0, 0, startOfNextWeek.Location())
					numDays = 7

				default:
					newStart = m.start.AddDate(0, 0, 1*(m.numDays))
					numDays = m.numDays
				}
				newEnd = newStart.AddDate(0, 0, numDays)
				cmds = append(cmds, getRecordsData(m.typ, m.db, m.period, newStart, newEnd, numDays, m.plain))
				m.busy = true
			}
		case "ctrl+t":
			if !m.busy {
				var start, end time.Time
				var numDays int

				switch m.period {
				case types.TimePeriodWeek:
					now := time.Now()
					weekday := now.Weekday()
					offset := (7 + weekday - time.Monday) % 7
					startOfWeek := now.AddDate(0, 0, -int(offset))
					start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
					numDays = 7
				default:
					now := time.Now()
					nDaysBack := now.AddDate(0, 0, -1*(m.numDays-1))

					start = time.Date(nDaysBack.Year(), nDaysBack.Month(), nDaysBack.Day(), 0, 0, 0, 0, nDaysBack.Location())
					numDays = m.numDays
				}
				end = start.AddDate(0, 0, numDays)
				cmds = append(cmds, getRecordsData(m.typ, m.db, m.period, start, end, numDays, m.plain))
				m.busy = true
			}
		}
	case recordsDataFetchedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.quitting = true
			return m, tea.Quit
		} else {
			m.start = msg.start
			m.end = msg.end
			m.report = msg.report
			m.busy = false
		}
	}
	return m, tea.Batch(cmds...)
}
