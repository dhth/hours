package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/hours/internal/types"
)

const (
	ctrlC                 = "ctrl+c"
	enter                 = "enter"
	escape                = "esc"
	viewPortMoveLineCount = 3
	msgCouldntSelectATask = "Couldn't select a task"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.frameCounter++
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// early check for window resizing and handling insufficient dimensions
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResizing(msg)
	case tea.KeyMsg:
		if msg.String() == ctrlC {
			return m, tea.Quit
		}

		if m.activeView == insufficientDimensionsView {
			switch msg.String() {
			case "q", escape:
				return m, tea.Quit
			default:
				return m, tea.Batch(cmds...)
			}
		}
	}

	if m.activeView != insufficientDimensionsView {
		if m.message.framesLeft > 0 {
			m.message.framesLeft--
		}

		if m.message.framesLeft == 0 {
			m.message.value = ""
		}
	}

	keyMsg, keyMsgOK := msg.(tea.KeyMsg)
	if keyMsgOK {
		if m.activeTasksList.FilterState() == list.Filtering {
			m.activeTasksList, cmd = m.activeTasksList.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch keyMsg.String() {
		case enter, "ctrl+s":
			var bail bool
			if keyMsg.String() == enter {
				switch m.activeView {
				case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
					if m.trackingFocussedField == entryComment {
						bail = true
					}
				}
			}

			if bail {
				break
			}

			var updateCmd tea.Cmd
			switch m.activeView {
			case taskInputView:
				updateCmd = m.getCmdToCreateOrUpdateTask()
			case editActiveTLView:
				updateCmd = m.getCmdToUpdateActiveTL()
			case finishActiveTLView:
				updateCmd = m.getCmdToFinishTrackingActiveTL()
			case manualTasklogEntryView, editSavedTLView:
				updateCmd = m.getCmdToCreateOrEditTL()
			}
			if updateCmd != nil {
				cmds = append(cmds, updateCmd)
				return m, tea.Batch(cmds...)
			}
		case escape:
			switch m.activeView {
			case taskInputView, editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				m.handleEscapeInForms()
				return m, tea.Batch(cmds...)
			}
		case "tab":
			m.goForwardInView()
		case "shift+tab":
			m.goBackwardInView()
		case "k":
			switch m.activeView {
			case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				err := m.shiftTime(types.ShiftBackward, types.ShiftMinute)
				if err != nil {
					return m, tea.Batch(cmds...)
				}
			}
		case "j":
			switch m.activeView {
			case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				err := m.shiftTime(types.ShiftForward, types.ShiftMinute)
				if err != nil {
					return m, tea.Batch(cmds...)
				}
			}
		case "K":
			switch m.activeView {
			case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				err := m.shiftTime(types.ShiftBackward, types.ShiftFiveMinutes)
				if err != nil {
					return m, tea.Batch(cmds...)
				}
			}
		case "J":
			switch m.activeView {
			case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				err := m.shiftTime(types.ShiftForward, types.ShiftFiveMinutes)
				if err != nil {
					return m, tea.Batch(cmds...)
				}
			}
		case "h":
			switch m.activeView {
			case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				err := m.shiftTime(types.ShiftBackward, types.ShiftDay)
				if err != nil {
					return m, tea.Batch(cmds...)
				}
			case taskLogDetailsView:
				m.taskLogList.CursorUp()
				m.handleRequestToViewTLDetails()
			}
		case "l":
			switch m.activeView {
			case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
				err := m.shiftTime(types.ShiftForward, types.ShiftDay)
				if err != nil {
					return m, tea.Batch(cmds...)
				}
			case taskLogDetailsView:
				m.taskLogList.CursorDown()
				m.handleRequestToViewTLDetails()
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
	case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
		for i := range m.tLInputs {
			m.tLInputs[i], cmd = m.tLInputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		m.tLCommentInput, cmd = m.tLCommentInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", escape:
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
		case "f":
			if m.activeView != taskListView {
				break
			}

			if !m.trackingActive {
				m.message = errMsg("Nothing is being tracked right now")
				break
			}

			handleCmd := m.getCmdToFinishActiveTLWithoutComment()
			if handleCmd != nil {
				cmds = append(cmds, handleCmd)
			}
		case "ctrl+s":
			switch m.activeView {
			case taskListView:
				switch m.trackingActive {
				case true:
					m.handleRequestToEditActiveTL()
				case false:
					m.handleRequestToCreateManualTL()
				}
			case taskLogView:
				m.handleRequestToEditSavedTL()
			}
		case "u":
			switch m.activeView {
			case taskListView:
				m.handleRequestToUpdateTask()
			case taskLogView:
				m.handleRequestToEditSavedTL()
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
		case "S":
			if m.activeView != taskListView {
				break
			}
			quickSwitchCmd := m.getCmdToQuickSwitchTracking()
			if quickSwitchCmd != nil {
				cmds = append(cmds, quickSwitchCmd)
			}
		case "a":
			if m.activeView == taskListView {
				m.handleRequestToCreateTask()
			}
		case "c":
			if m.activeView == taskListView || m.activeView == inactiveTaskListView {
				m.handleCopyTaskSummary()
			}
		case "k":
			m.handleRequestToScrollVPUp()
		case "j":
			m.handleRequestToScrollVPDown()
		case "d":
			if m.activeView == taskLogView {
				m.handleRequestToViewTLDetails()
			}
		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
		}
	case taskCreatedMsg:
		if msg.err != nil {
			m.message = errMsg(fmt.Sprintf("Error creating task: %s", msg.err))
		} else {
			cmds = append(cmds, fetchTasks(m.db, true))
		}
	case taskUpdatedMsg:
		if msg.err != nil {
			m.message = errMsg(fmt.Sprintf("Error updating task: %s", msg.err))
		} else {
			msg.tsk.Summary = msg.summary
			msg.tsk.UpdateListTitle()
		}
	case tasksFetchedMsg:
		handleCmd := m.handleTasksFetchedMsg(msg)
		if handleCmd != nil {
			cmds = append(cmds, handleCmd)
		}
	case activeTLUpdatedMsg:
		if msg.err != nil {
			m.message = errMsg(msg.err.Error())
		} else {
			m.activeTLBeginTS = msg.beginTS
			m.activeTLComment = msg.comment
		}
	case manualTLInsertedMsg:
		handleCmds := m.handleManualTLInsertedMsg(msg)
		if handleCmds != nil {
			cmds = append(cmds, handleCmds...)
		}
	case savedTLEditedMsg:
		handleCmds := m.handleSavedTLEditedMsg(msg)
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
	case activeTLSwitchedMsg:
		updateCmd := m.handleActiveTLSwitchedMsg(msg)
		if updateCmd != nil {
			cmds = append(cmds, updateCmd)
		}
	case taskRepUpdatedMsg:
		if msg.err != nil {
			m.message = errMsg(fmt.Sprintf("Error updating task status: %s", msg.err))
		} else {
			msg.tsk.UpdateListDesc(m.timeProvider)
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
			m.message = errMsg("Error updating task's active status: " + msg.err.Error())
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
		case ctrlC, "q", escape:
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			if !m.busy {
				var dr types.DateRange

				switch m.period {
				case types.TimePeriodWeek:
					weekday := m.dateRange.Start.Weekday()
					offset := (7 + weekday - time.Monday) % 7
					startOfPrevWeek := m.dateRange.Start.AddDate(0, 0, -int(offset+7))
					dr.Start = time.Date(startOfPrevWeek.Year(), startOfPrevWeek.Month(), startOfPrevWeek.Day(), 0, 0, 0, 0, startOfPrevWeek.Location())
				default:
					dr.Start = m.dateRange.Start.AddDate(0, 0, -m.dateRange.NumDays)
				}

				dr.NumDays = m.dateRange.NumDays
				dr.End = dr.Start.AddDate(0, 0, m.dateRange.NumDays)
				cmds = append(cmds, getRecordsData(m.kind, m.db, m.style, dr, m.taskStatus, m.plain))
				m.busy = true
			}
		case "right", "l":
			if !m.busy {
				var dr types.DateRange

				switch m.period {
				case types.TimePeriodWeek:
					weekday := m.dateRange.Start.Weekday()
					offset := (7 + weekday - time.Monday) % 7
					startOfNextWeek := m.dateRange.Start.AddDate(0, 0, 7-int(offset))
					dr.Start = time.Date(startOfNextWeek.Year(), startOfNextWeek.Month(), startOfNextWeek.Day(), 0, 0, 0, 0, startOfNextWeek.Location())
					dr.NumDays = 7

				default:
					dr.Start = m.dateRange.Start.AddDate(0, 0, 1*(m.dateRange.NumDays))
				}

				dr.NumDays = m.dateRange.NumDays
				dr.End = dr.Start.AddDate(0, 0, dr.NumDays)
				cmds = append(cmds, getRecordsData(m.kind, m.db, m.style, dr, m.taskStatus, m.plain))
				m.busy = true
			}
		case "ctrl+t":
			if !m.busy {
				var dr types.DateRange

				now := m.timeProvider.Now()
				switch m.period {
				case types.TimePeriodWeek:
					weekday := now.Weekday()
					offset := (7 + weekday - time.Monday) % 7
					startOfWeek := now.AddDate(0, 0, -int(offset))
					dr.Start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
					dr.NumDays = 7
				default:
					nDaysBack := now.AddDate(0, 0, -1*(m.dateRange.NumDays-1))

					dr.Start = time.Date(nDaysBack.Year(), nDaysBack.Month(), nDaysBack.Day(), 0, 0, 0, 0, nDaysBack.Location())
				}

				dr.NumDays = m.dateRange.NumDays
				dr.End = dr.Start.AddDate(0, 0, dr.NumDays)
				cmds = append(cmds, getRecordsData(m.kind, m.db, m.style, dr, m.taskStatus, m.plain))
				m.busy = true
			}
		}
	case recordsDataFetchedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.quitting = true
			return m, tea.Quit
		}

		m.dateRange = msg.dateRange
		m.report = msg.report
		m.busy = false
	}
	return m, tea.Batch(cmds...)
}
