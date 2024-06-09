package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

const useHighPerformanceRenderer = false

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.message = ""
	m.errorMessage = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.taskList.FilterState() == list.Filtering {
			m.taskList, cmd = m.taskList.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch m.activeView {
			case taskInputView:
				m.activeView = taskListView
				if m.taskInputs[summaryField].Value() != "" {
					switch m.taskMgmtContext {
					case taskCreateCxt:
						cmds = append(cmds, createTask(m.db, m.taskInputs[summaryField].Value()))
						m.taskInputs[summaryField].SetValue("")
					case taskUpdateCxt:
						selectedTask, ok := m.taskList.SelectedItem().(*task)
						if ok {
							cmds = append(cmds, updateTask(m.db, selectedTask, m.taskInputs[summaryField].Value()))
							m.taskInputs[summaryField].SetValue("")
						}
					}
					return m, tea.Batch(cmds...)
				}
			case askForCommentView:
				m.activeView = taskListView
				if m.trackingInputs[entryComment].Value() != "" {
					m.activeTLEndTS = time.Now()
					cmds = append(cmds, toggleTracking(m.db, m.activeTaskId, m.activeTLBeginTS, m.activeTLEndTS, m.trackingInputs[entryComment].Value()))
					m.trackingInputs[entryComment].SetValue("")
					return m, tea.Batch(cmds...)
				}
			case manualTasklogEntryView:
				beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
				if err != nil {
					m.errorMessage = err.Error()
					return m, tea.Batch(cmds...)
				}

				endTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryEndTS].Value(), time.Local)

				if err != nil {
					m.errorMessage = err.Error()
					return m, tea.Batch(cmds...)
				}

				comment := m.trackingInputs[entryComment].Value()

				if len(comment) == 0 {
					m.errorMessage = "Comment cannot be empty"
					return m, tea.Batch(cmds...)

				}

				for i := range m.trackingInputs {
					m.trackingInputs[i].SetValue("")
				}
				task, ok := m.taskList.SelectedItem().(*task)
				if ok {
					switch m.tasklogSaveType {
					case tasklogInsert:
						cmds = append(cmds, insertManualEntry(m.db, task.id, beginTS.Local(), endTS.Local(), comment))
						m.activeView = taskListView
					}
				}
				return m, tea.Batch(cmds...)
			}
		case "esc":
			switch m.activeView {
			case taskInputView:
				m.activeView = taskListView
				for i := range m.taskInputs {
					m.taskInputs[i].SetValue("")
				}
			case askForCommentView:
				m.activeView = taskListView
				m.trackingInputs[entryComment].SetValue("")
			case manualTasklogEntryView:
				switch m.tasklogSaveType {
				case tasklogInsert:
					m.activeView = taskListView
				}
				for i := range m.trackingInputs {
					m.trackingInputs[i].SetValue("")
				}
			}
		case "tab":
			switch m.activeView {
			case taskListView:
				m.activeView = taskLogView
				cmds = append(cmds, fetchTaskLogEntries(m.db))
			case taskLogView:
				m.activeView = taskListView
			case manualTasklogEntryView:
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
				m.activeView = taskListView
			case taskListView:
				m.activeView = taskLogView
				cmds = append(cmds, fetchTaskLogEntries(m.db))
			case manualTasklogEntryView:
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
		}
	}

	switch m.activeView {
	case taskInputView:
		for i := range m.taskInputs {
			m.taskInputs[i], cmd = m.taskInputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case askForCommentView:
		m.trackingInputs[entryComment], cmd = m.trackingInputs[entryComment].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case manualTasklogEntryView:
		for i := range m.trackingInputs {
			m.trackingInputs[i], cmd = m.trackingInputs[i].Update(msg)
		}
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			switch m.activeView {
			case taskListView:
				fs := m.taskList.FilterState()
				if fs == list.Filtering || fs == list.FilterApplied {
					m.taskList.ResetFilter()
				} else {
					return m, tea.Quit
				}
			case taskLogView:
				fs := m.taskLogList.FilterState()
				if fs == list.Filtering || fs == list.FilterApplied {
					m.taskLogList.ResetFilter()
				} else {
					return m, tea.Quit
				}
			case helpView:
				m.activeView = taskListView
			default:
				return m, tea.Quit
			}
		case "1":
			if m.activeView != taskListView {
				m.activeView = taskListView
			}
		case "2":
			if m.activeView != taskLogView {
				m.activeView = taskLogView
				cmds = append(cmds, fetchTaskLogEntries(m.db))
			}
		case "ctrl+r":
			switch m.activeView {
			case taskListView:
				cmds = append(cmds, fetchTasks(m.db))
			case taskLogView:
				cmds = append(cmds, fetchTaskLogEntries(m.db))
				m.taskLogList.ResetSelected()
			}
		case "ctrl+t":
			if m.activeView == taskListView {
				if m.trackingActive {
					if m.taskList.IsFiltered() {
						m.taskList.ResetFilter()
					}
					activeIndex, ok := m.taskIndexMap[m.activeTaskId]
					if ok {
						m.taskList.Select(activeIndex)
					}
				} else {
					m.message = "Nothing is being tracked right now"
				}
			}
		case "ctrl+s":
			if m.activeView == taskListView && !m.trackingActive {
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
		case "ctrl+d":
			switch m.activeView {
			case taskLogView:
				entry, ok := m.taskLogList.SelectedItem().(taskLogEntry)
				if ok {
					cmds = append(cmds, deleteLogEntry(m.db, &entry))
					return m, tea.Batch(cmds...)
				} else {
					msg := "Couldn't delete task log entry"
					m.message = msg
					m.messages = append(m.messages, msg)
				}
			}
		case "s":
			switch m.activeView {
			case taskListView:
				if m.taskList.FilterState() != list.Filtering {
					if m.changesLocked {
						message := "Changes locked momentarily"
						m.message = message
						m.messages = append(m.messages, message)
					}
					task, ok := m.taskList.SelectedItem().(*task)
					if !ok {
						message := "Couldn't select a task"
						m.message = message
						m.messages = append(m.messages, message)
					} else {
						if m.lastChange == updateChange {
							m.changesLocked = true
							m.activeTLBeginTS = time.Now()
							cmds = append(cmds, toggleTracking(m.db, task.id, m.activeTLBeginTS, m.activeTLEndTS, ""))
						} else if m.lastChange == insertChange {
							m.activeView = askForCommentView
							m.trackingFocussedField = entryComment
							m.trackingInputs[m.trackingFocussedField].Focus()
						}
					}
				}
			}
		case "a":
			switch m.activeView {
			case taskListView:
				if m.taskList.FilterState() != list.Filtering {
					if m.changesLocked {
						message := "Changes locked momentarily"
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
			switch m.activeView {
			case taskListView:
				if m.taskList.FilterState() != list.Filtering {
					if m.changesLocked {
						message := "Changes locked momentarily"
						m.message = message
						m.messages = append(m.messages, message)
					}
					task, ok := m.taskList.SelectedItem().(*task)
					if ok {
						m.activeView = taskInputView
						m.taskInputFocussedField = summaryField
						m.taskInputs[summaryField].Focus()
						m.taskInputs[summaryField].SetValue(task.summary)
						m.taskMgmtContext = taskUpdateCxt
					} else {
						m.message = "Couldn't select a task"
					}
				}
			}
		case "?":
			m.lastView = m.activeView
			m.activeView = helpView
		}

	case tea.WindowSizeMsg:
		w, h := stackListStyle.GetFrameSize()
		m.terminalHeight = msg.Height
		m.taskList.SetHeight(msg.Height - h - 2)
		m.taskLogList.SetHeight(msg.Height - h - 2)
		m.taskLogList.SetHeight(msg.Height - h - 2)
		if !m.helpVPReady {
			m.helpVP = viewport.New(w-5, m.terminalHeight-7)
			m.helpVP.HighPerformanceRendering = useHighPerformanceRenderer
			m.helpVP.SetContent(helpText)
			m.helpVPReady = true
		} else {
			m.helpVP.Height = m.terminalHeight - 7
			m.helpVP.Width = w - 5

		}
	case taskCreatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error creating task: %s", msg.err)
		} else {
			cmds = append(cmds, fetchTasks(m.db))
		}
	case taskUpdatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error updating task: %s", msg.err)
		} else {
			msg.tsk.summary = msg.summary
			msg.tsk.updateTitle()
		}
	case tasksFetched:
		if msg.err != nil {
			message := "error fetching tasks : " + msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
		} else {
			m.taskMap = make(map[int]*task)
			m.taskIndexMap = make(map[int]int)
			tasks := make([]list.Item, len(msg.tasks))
			for i, task := range msg.tasks {
				task.updateTitle()
				task.updateDesc()
				tasks[i] = &task
				m.taskMap[task.id] = &task
				m.taskIndexMap[task.id] = i
			}
			m.taskList.SetItems(tasks)
			m.taskList.Title = "Tasks"
			m.tasksFetched = true

			cmds = append(cmds, fetchActiveTask(m.db))
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
			task, ok := m.taskMap[msg.taskId]

			if ok {
				cmds = append(cmds, updateTaskRep(m.db, task))
			}
		}
	case taskLogEntriesFetchedMsg:
		if msg.err != nil {
			message := msg.err.Error()
			m.message = "Error fetching synced task log entries: " + message
			m.messages = append(m.messages, message)
		} else {
			var items []list.Item
			for _, e := range msg.entries {
				items = append(items, list.Item(e))
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
				m.activeTaskId = msg.activeTaskId
				m.lastChange = insertChange
				m.activeTLBeginTS = msg.beginTs
				activeTask, ok := m.taskMap[m.activeTaskId]
				if ok {
					activeTask.trackingActive = true
					activeTask.updateTitle()

					// go to tracked item on startup
					activeIndex, ok := m.taskIndexMap[msg.activeTaskId]
					if ok {
						m.taskList.Select(activeIndex)
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

			task, ok := m.taskMap[msg.taskId]

			if ok {
				if msg.finished {
					m.lastChange = updateChange
					task.trackingActive = false
					m.trackingActive = false
					m.activeTaskId = -1
					cmds = append(cmds, updateTaskRep(m.db, task))
				} else {
					m.lastChange = insertChange
					task.trackingActive = true
					m.trackingActive = true
					m.activeTaskId = msg.taskId
				}
				task.updateTitle()
			}
		}
	case taskRepUpdatedMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("Error updating task status: %s", msg.err)
		} else {
			msg.tsk.updateDesc()
		}
	case taskLogEntryDeletedMsg:
		if msg.err != nil {
			message := "error deleting entry: " + msg.err.Error()
			m.message = message
			m.messages = append(m.messages, message)
		} else {
			task, ok := m.taskMap[msg.entry.taskId]
			if ok {
				cmds = append(cmds, updateTaskRep(m.db, task))
			}
			cmds = append(cmds, fetchTaskLogEntries(m.db))
		}
	case HideHelpMsg:
		m.showHelpIndicator = false
	}

	switch m.activeView {
	case taskListView:
		m.taskList, cmd = m.taskList.Update(msg)
		cmds = append(cmds, cmd)
	case taskLogView:
		m.taskLogList, cmd = m.taskLogList.Update(msg)
		cmds = append(cmds, cmd)
	case helpView:
		m.helpVP, cmd = m.helpVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
