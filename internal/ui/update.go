package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
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
			switch m.activeView {
			case taskInputView:
				m.activeView = activeTaskListView
				if m.taskInputs[summaryField].Value() != "" {
					switch m.taskMgmtContext {
					case taskCreateCxt:
						cmds = append(cmds, createTask(m.db, m.taskInputs[summaryField].Value()))
						m.taskInputs[summaryField].SetValue("")
					case taskUpdateCxt:
						selectedTask, ok := m.activeTasksList.SelectedItem().(*types.Task)
						if ok {
							cmds = append(cmds, updateTask(m.db, selectedTask, m.taskInputs[summaryField].Value()))
							m.taskInputs[summaryField].SetValue("")
						}
					}
					return m, tea.Batch(cmds...)
				}
			case editStartTsView:
				beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
				if err != nil {
					m.message = err.Error()
					return m, tea.Batch(cmds...)
				}

				cmds = append(cmds, updateTLBeginTS(m.db, beginTS))
				m.trackingInputs[entryBeginTS].SetValue("")
				m.activeView = activeTaskListView

				return m, tea.Batch(cmds...)
			case askForCommentView:
				beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
				if err != nil {
					m.message = err.Error()
					return m, tea.Batch(cmds...)
				}
				m.activeTLBeginTS = beginTS

				endTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryEndTS].Value(), time.Local)
				if err != nil {
					m.message = err.Error()
					return m, tea.Batch(cmds...)
				}
				m.activeTLEndTS = endTS

				if m.activeTLEndTS.Sub(m.activeTLBeginTS).Seconds() <= 0 {
					m.message = "time spent needs to be positive"
					return m, tea.Batch(cmds...)
				}

				if m.trackingInputs[entryComment].Value() == "" {
					m.message = "Comment cannot be empty"
					return m, tea.Batch(cmds...)
				}

				cmds = append(cmds, toggleTracking(m.db, m.activeTaskID, m.activeTLBeginTS, m.activeTLEndTS, m.trackingInputs[entryComment].Value()))
				m.activeView = activeTaskListView

				for i := range m.trackingInputs {
					m.trackingInputs[i].SetValue("")
				}
				return m, tea.Batch(cmds...)

			case manualTasklogEntryView:
				beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
				if err != nil {
					m.message = err.Error()
					return m, tea.Batch(cmds...)
				}

				endTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryEndTS].Value(), time.Local)
				if err != nil {
					m.message = err.Error()
					return m, tea.Batch(cmds...)
				}

				if endTS.Sub(beginTS).Seconds() <= 0 {
					m.message = "time spent needs to be positive"
					return m, tea.Batch(cmds...)
				}

				comment := m.trackingInputs[entryComment].Value()

				if len(comment) == 0 {
					m.message = "Comment cannot be empty"
					return m, tea.Batch(cmds...)
				}

				task, ok := m.activeTasksList.SelectedItem().(*types.Task)
				if ok && m.tasklogSaveType == tasklogInsert {
					cmds = append(cmds, insertManualEntry(m.db, task.ID, beginTS, endTS, comment))
					m.activeView = activeTaskListView
				}
				for i := range m.trackingInputs {
					m.trackingInputs[i].SetValue("")
				}
				return m, tea.Batch(cmds...)
			}
		case "esc":
			switch m.activeView {
			case taskInputView:
				m.activeView = activeTaskListView
				for i := range m.taskInputs {
					m.taskInputs[i].SetValue("")
				}
			case editStartTsView:
				m.taskInputs[entryBeginTS].SetValue("")
				m.activeView = activeTaskListView
			case askForCommentView:
				m.activeView = activeTaskListView
				m.trackingInputs[entryComment].SetValue("")
			case manualTasklogEntryView:
				if m.tasklogSaveType == tasklogInsert {
					m.activeView = activeTaskListView
				}
				for i := range m.trackingInputs {
					m.trackingInputs[i].SetValue("")
				}
			}
		case "tab":
			switch m.activeView {
			case activeTaskListView:
				m.activeView = taskLogView
			case taskLogView:
				m.activeView = inactiveTaskListView
			case inactiveTaskListView:
				m.activeView = activeTaskListView
			case askForCommentView, manualTasklogEntryView:
				switch m.trackingFocussedField {
				case entryBeginTS:
					m.trackingFocussedField = entryEndTS
				case entryEndTS:
					m.trackingFocussedField = entryComment
				case entryComment:
					m.trackingFocussedField = entryBeginTS
				}
				for i := range m.trackingInputs {
					m.trackingInputs[i].Blur()
				}
				m.trackingInputs[m.trackingFocussedField].Focus()
			}
		case "shift+tab":
			switch m.activeView {
			case taskLogView:
				m.activeView = activeTaskListView
			case activeTaskListView:
				m.activeView = inactiveTaskListView
			case inactiveTaskListView:
				m.activeView = taskLogView
			case askForCommentView, manualTasklogEntryView:
				switch m.trackingFocussedField {
				case entryBeginTS:
					m.trackingFocussedField = entryComment
				case entryEndTS:
					m.trackingFocussedField = entryBeginTS
				case entryComment:
					m.trackingFocussedField = entryEndTS
				}
				for i := range m.trackingInputs {
					m.trackingInputs[i].Blur()
				}
				m.trackingInputs[m.trackingFocussedField].Focus()
			}
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
		}
	}

	switch m.activeView {
	case taskInputView:
		for i := range m.taskInputs {
			m.taskInputs[i], cmd = m.taskInputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case editStartTsView:
		m.trackingInputs[entryBeginTS], cmd = m.trackingInputs[entryBeginTS].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case askForCommentView:
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
			switch m.activeView {
			case activeTaskListView:
				fs := m.activeTasksList.FilterState()
				if fs == list.Filtering || fs == list.FilterApplied {
					m.activeTasksList.ResetFilter()
				} else {
					return m, tea.Quit
				}
			case taskLogView:
				fs := m.taskLogList.FilterState()
				if fs == list.Filtering || fs == list.FilterApplied {
					m.taskLogList.ResetFilter()
				} else {
					m.activeView = activeTaskListView
				}
			case inactiveTaskListView:
				fs := m.inactiveTasksList.FilterState()
				if fs == list.Filtering || fs == list.FilterApplied {
					m.inactiveTasksList.ResetFilter()
				} else {
					m.activeView = activeTaskListView
				}
			case helpView:
				m.activeView = activeTaskListView
			default:
				return m, tea.Quit
			}
		case "1":
			if m.activeView != activeTaskListView {
				m.activeView = activeTaskListView
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
			switch m.activeView {
			case activeTaskListView:
				cmds = append(cmds, fetchTasks(m.db, true))
			case taskLogView:
				cmds = append(cmds, fetchTaskLogEntries(m.db))
				m.taskLogList.ResetSelected()
			case inactiveTaskListView:
				cmds = append(cmds, fetchTasks(m.db, false))
				m.inactiveTasksList.ResetSelected()
			}
		case "ctrl+t":
			if m.activeView == activeTaskListView {
				if m.trackingActive {
					if m.activeTasksList.IsFiltered() {
						m.activeTasksList.ResetFilter()
					}
					activeIndex, ok := m.activeTaskIndexMap[m.activeTaskID]
					if ok {
						m.activeTasksList.Select(activeIndex)
					}
				} else {
					m.message = "Nothing is being tracked right now"
				}
			}
		case "ctrl+s":
			if m.activeView == activeTaskListView {
				_, ok := m.activeTasksList.SelectedItem().(*types.Task)
				if !ok {
					message := msgCouldntSelectATask
					m.message = message
					m.messages = append(m.messages, message)
				} else {
					if m.trackingActive {
						m.activeView = editStartTsView
						m.trackingFocussedField = entryBeginTS
						m.trackingInputs[entryBeginTS].SetValue(m.activeTLBeginTS.Format(timeFormat))
						m.trackingInputs[m.trackingFocussedField].Focus()
					} else {
						m.activeView = manualTasklogEntryView
						m.tasklogSaveType = tasklogInsert
						m.trackingFocussedField = entryBeginTS
						currentTime := time.Now()
						dateString := currentTime.Format("2006/01/02")
						currentTimeStr := currentTime.Format(timeFormat)

						m.trackingInputs[entryBeginTS].SetValue(dateString + " ")
						m.trackingInputs[entryEndTS].SetValue(currentTimeStr)

						for i := range m.trackingInputs {
							m.trackingInputs[i].Blur()
						}
						m.trackingInputs[m.trackingFocussedField].Focus()
					}
				}
			}
		case "ctrl+d":
			switch m.activeView {
			case activeTaskListView:
				task, ok := m.activeTasksList.SelectedItem().(*types.Task)
				if ok {
					if task.TrackingActive {
						m.message = "Cannot deactivate a task being tracked; stop tracking and try again."
					} else {
						cmds = append(cmds, updateTaskActiveStatus(m.db, task, false))
					}
				} else {
					msg := msgCouldntSelectATask
					m.message = msg
					m.messages = append(m.messages, msg)
				}
			case taskLogView:
				entry, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
				if ok {
					cmds = append(cmds, deleteLogEntry(m.db, &entry))
				} else {
					msg := "Couldn't delete task log entry"
					m.message = msg
					m.messages = append(m.messages, msg)
				}
			case inactiveTaskListView:
				task, ok := m.inactiveTasksList.SelectedItem().(*types.Task)
				if ok {
					cmds = append(cmds, updateTaskActiveStatus(m.db, task, true))
				} else {
					msg := msgCouldntSelectATask
					m.message = msg
					m.messages = append(m.messages, msg)
				}
			}
		case "ctrl+x":
			if m.activeView == activeTaskListView && m.trackingActive {
				cmds = append(cmds, deleteActiveTaskLog(m.db))
			}
		case "s":
			if m.activeView == activeTaskListView {
				if m.activeTasksList.FilterState() != list.Filtering {
					if m.changesLocked {
						message := msgChangesLocked
						m.message = message
						m.messages = append(m.messages, message)
					}
					task, ok := m.activeTasksList.SelectedItem().(*types.Task)
					if !ok {
						message := "Couldn't select a task"
						m.message = message
						m.messages = append(m.messages, message)
					} else {
						if m.lastChange == updateChange {
							m.changesLocked = true
							m.activeTLBeginTS = time.Now()
							cmds = append(cmds, toggleTracking(m.db, task.ID, m.activeTLBeginTS, m.activeTLEndTS, ""))
						} else if m.lastChange == insertChange {
							m.activeView = askForCommentView
							m.activeTLEndTS = time.Now()

							beginTimeStr := m.activeTLBeginTS.Format(timeFormat)
							currentTimeStr := m.activeTLEndTS.Format(timeFormat)

							m.trackingInputs[entryBeginTS].SetValue(beginTimeStr)
							m.trackingInputs[entryEndTS].SetValue(currentTimeStr)
							m.trackingFocussedField = entryComment

							for i := range m.trackingInputs {
								m.trackingInputs[i].Blur()
							}
							m.trackingInputs[m.trackingFocussedField].Focus()
						}
					}
				}
			}
		case "a":
			if m.activeView == activeTaskListView {
				if m.activeTasksList.FilterState() != list.Filtering {
					if m.changesLocked {
						message := msgChangesLocked
						m.message = message
						m.messages = append(m.messages, message)
					}
					m.activeView = taskInputView
					m.taskInputFocussedField = summaryField
					m.taskInputs[summaryField].Focus()
					m.taskMgmtContext = taskCreateCxt
				}
			}
		case "u":
			if m.activeView == activeTaskListView {
				if m.activeTasksList.FilterState() != list.Filtering {
					if m.changesLocked {
						message := msgChangesLocked
						m.message = message
						m.messages = append(m.messages, message)
					}
					task, ok := m.activeTasksList.SelectedItem().(*types.Task)
					if ok {
						m.activeView = taskInputView
						m.taskInputFocussedField = summaryField
						m.taskInputs[summaryField].Focus()
						m.taskInputs[summaryField].SetValue(task.Summary)
						m.taskMgmtContext = taskUpdateCxt
					} else {
						m.message = "Couldn't select a task"
					}
				}
			}

		case "k":
			if m.activeView != helpView {
				break
			}
			if m.helpVP.AtTop() {
				break
			}
			m.helpVP.LineUp(viewPortMoveLineCount)

		case "j":
			if m.activeView != helpView {
				break
			}
			if m.helpVP.AtBottom() {
				break
			}
			m.helpVP.LineDown(viewPortMoveLineCount)

		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
		}

	case tea.WindowSizeMsg:
		w, h := listStyle.GetFrameSize()
		m.terminalHeight = msg.Height

		m.taskLogList.SetWidth(msg.Width - w)
		m.taskLogList.SetHeight(msg.Height - h - 2)

		m.activeTasksList.SetWidth(msg.Width - w)
		m.activeTasksList.SetHeight(msg.Height - h - 2)

		m.inactiveTasksList.SetWidth(msg.Width - w)
		m.inactiveTasksList.SetHeight(msg.Height - h - 2)

		if !m.helpVPReady {
			m.helpVP = viewport.New(w-5, m.terminalHeight-7)
			m.helpVP.SetContent(helpText)
			m.helpVP.KeyMap.Up.SetEnabled(false)
			m.helpVP.KeyMap.Down.SetEnabled(false)
			m.helpVPReady = true
		} else {
			m.helpVP.Height = m.terminalHeight - 7
			m.helpVP.Width = w - 5

		}
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
	case tasksFetched:
		if msg.err != nil {
			message := "error fetching tasks : " + msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
		} else {
			if msg.active {
				m.activeTaskMap = make(map[int]*types.Task)
				m.activeTaskIndexMap = make(map[int]int)
				tasks := make([]list.Item, len(msg.tasks))
				for i, task := range msg.tasks {
					task.UpdateTitle()
					task.UpdateDesc()
					tasks[i] = &task
					m.activeTaskMap[task.ID] = &task
					m.activeTaskIndexMap[task.ID] = i
				}
				m.activeTasksList.SetItems(tasks)
				m.activeTasksList.Title = "Tasks"
				m.tasksFetched = true
				cmds = append(cmds, fetchActiveTask(m.db))
			} else {
				inactiveTasks := make([]list.Item, len(msg.tasks))
				for i, inactiveTask := range msg.tasks {
					inactiveTask.UpdateTitle()
					inactiveTask.UpdateDesc()
					inactiveTasks[i] = &inactiveTask
				}
				m.inactiveTasksList.SetItems(inactiveTasks)
			}
		}
	case tlBeginTSUpdatedMsg:
		if msg.err != nil {
			message := msg.err.Error()
			m.message = "Error updating begin time: " + message
			m.messages = append(m.messages, message)
		} else {
			m.activeTLBeginTS = msg.beginTS
		}
	case manualTaskLogInserted:
		if msg.err != nil {
			message := msg.err.Error()
			m.message = "Error inserting task log: " + message
			m.messages = append(m.messages, message)
		} else {
			for i := range m.trackingInputs {
				m.trackingInputs[i].SetValue("")
			}
			task, ok := m.activeTaskMap[msg.taskID]

			if ok {
				cmds = append(cmds, updateTaskRep(m.db, task))
			}
			cmds = append(cmds, fetchTaskLogEntries(m.db))
		}
	case taskLogEntriesFetchedMsg:
		if msg.err != nil {
			message := msg.err.Error()
			m.message = "Error fetching task log entries: " + message
			m.messages = append(m.messages, message)
		} else {
			var items []list.Item
			for _, e := range msg.entries {
				e.UpdateTitle()
				e.UpdateDesc()
				items = append(items, e)
			}
			m.taskLogList.SetItems(items)
		}
	case activeTaskFetchedMsg:
		if msg.err != nil {
			message := msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
		} else {
			if msg.noneActive {
				m.lastChange = updateChange
			} else {
				m.activeTaskID = msg.activeTaskID
				m.lastChange = insertChange
				m.activeTLBeginTS = msg.beginTs
				activeTask, ok := m.activeTaskMap[m.activeTaskID]
				if ok {
					activeTask.TrackingActive = true
					activeTask.UpdateTitle()

					// go to tracked item on startup
					activeIndex, ok := m.activeTaskIndexMap[msg.activeTaskID]
					if ok {
						m.activeTasksList.Select(activeIndex)
					}
				}
				m.trackingActive = true
			}
		}
	case trackingToggledMsg:
		if msg.err != nil {
			message := msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
			m.trackingActive = false
		} else {
			m.changesLocked = false

			task, ok := m.activeTaskMap[msg.taskID]

			if ok {
				if msg.finished {
					m.lastChange = updateChange
					task.TrackingActive = false
					m.trackingActive = false
					m.activeTaskID = -1
					cmds = append(cmds, updateTaskRep(m.db, task))
					cmds = append(cmds, fetchTaskLogEntries(m.db))
				} else {
					m.lastChange = insertChange
					task.TrackingActive = true
					m.trackingActive = true
					m.activeTaskID = msg.taskID
				}
				task.UpdateTitle()
			}
		}
	case taskRepUpdatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error updating task status: %s", msg.err)
		} else {
			msg.tsk.UpdateDesc()
		}
	case taskLogEntryDeletedMsg:
		if msg.err != nil {
			message := "error deleting entry: " + msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
		} else {
			task, ok := m.activeTaskMap[msg.entry.TaskID]
			if ok {
				cmds = append(cmds, updateTaskRep(m.db, task))
			}
			cmds = append(cmds, fetchTaskLogEntries(m.db))
		}
	case activeTaskLogDeletedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error deleting active log entry: %s", msg.err)
		} else {
			activeTask, ok := m.activeTaskMap[m.activeTaskID]
			if ok {
				activeTask.TrackingActive = false
				activeTask.UpdateTitle()
			}
			m.lastChange = updateChange
			m.trackingActive = false
			m.activeTaskID = -1
		}
	case taskActiveStatusUpdated:
		if msg.err != nil {
			message := "error updating task's active status: " + msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
		} else {
			cmds = append(cmds, fetchTasks(m.db, true))
			cmds = append(cmds, fetchTasks(m.db, false))
		}
	case HideHelpMsg:
		m.showHelpIndicator = false
	}

	switch m.activeView {
	case activeTaskListView:
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

func (m Model) shiftTime(direction types.TimeShiftDirection, duration types.TimeShiftDuration) error {
	if m.activeView == editStartTsView || m.activeView == askForCommentView || m.activeView == manualTasklogEntryView {
		if m.trackingFocussedField == entryBeginTS || m.trackingFocussedField == entryEndTS {
			ts, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[m.trackingFocussedField].Value(), time.Local)
			if err != nil {
				return err
			}

			newTs := types.GetShiftedTime(ts, direction, duration)

			m.trackingInputs[m.trackingFocussedField].SetValue(newTs.Format(timeFormat))
		}
	}
	return nil
}
