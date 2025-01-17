package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/hours/internal/types"
)

const (
	genericErrorMsg = "Something went wrong"
	removeFilterMsg = "Remove filter first"
)

func (m *Model) getCmdToCreateOrUpdateTask() tea.Cmd {
	if strings.TrimSpace(m.taskInputs[summaryField].Value()) == "" {
		m.message = "Task summary cannot be empty"
		return nil
	}

	var cmd tea.Cmd
	switch m.taskMgmtContext {
	case taskCreateCxt:
		cmd = createTask(m.db, m.taskInputs[summaryField].Value())
		m.taskInputs[summaryField].SetValue("")
	case taskUpdateCxt:
		selectedTask, ok := m.activeTasksList.SelectedItem().(*types.Task)
		if !ok {
			m.message = "Something went wrong"
			return nil
		}
		cmd = updateTask(m.db, selectedTask, m.taskInputs[summaryField].Value())
		m.taskInputs[summaryField].SetValue("")
	}

	m.activeView = taskListView
	return cmd
}

func (m *Model) getCmdToUpdateActiveTL() tea.Cmd {
	beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}

	m.trackingInputs[entryBeginTS].SetValue("")
	m.activeView = taskListView
	return updateTLBeginTS(m.db, beginTS)
}

func (m *Model) getCmdToSaveActiveTL() tea.Cmd {
	beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}
	m.activeTLBeginTS = beginTS

	endTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryEndTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}
	m.activeTLEndTS = endTS

	if m.activeTLEndTS.Sub(m.activeTLBeginTS).Seconds() <= 0 {
		m.message = "time spent needs to be positive"
		return nil
	}

	if m.trackingInputs[entryComment].Value() == "" {
		m.message = "Comment cannot be empty"
		return nil
	}

	comment := m.trackingInputs[entryComment].Value()
	for i := range m.trackingInputs {
		m.trackingInputs[i].SetValue("")
	}
	m.activeView = taskListView
	return toggleTracking(m.db, m.activeTaskID, m.activeTLBeginTS, m.activeTLEndTS, comment)
}

func (m *Model) getCmdToSaveOrUpdateTL() tea.Cmd {
	beginTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryBeginTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}

	endTS, err := time.ParseInLocation(string(timeFormat), m.trackingInputs[entryEndTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}

	if endTS.Sub(beginTS).Seconds() <= 0 {
		m.message = "time spent needs to be positive"
		return nil
	}

	comment := m.trackingInputs[entryComment].Value()

	if len(comment) == 0 {
		m.message = "Comment cannot be empty"
		return nil
	}

	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = "Something went wrong"
		return nil
	}
	if m.tasklogSaveType != tasklogInsert {
		return nil
	}
	for i := range m.trackingInputs {
		m.trackingInputs[i].SetValue("")
	}

	m.activeView = taskListView
	return insertManualTL(m.db, task.ID, beginTS, endTS, comment)
}

func (m *Model) handleEscape() {
	switch m.activeView {
	case taskInputView:
		m.activeView = taskListView
		for i := range m.taskInputs {
			m.taskInputs[i].SetValue("")
		}
	case editActiveTLView:
		m.taskInputs[entryBeginTS].SetValue("")
		m.activeView = taskListView
	case saveActiveTLView:
		m.activeView = taskListView
		m.trackingInputs[entryComment].SetValue("")
	case manualTasklogEntryView:
		if m.tasklogSaveType == tasklogInsert {
			m.activeView = taskListView
		}
		for i := range m.trackingInputs {
			m.trackingInputs[i].SetValue("")
		}
	}
}

func (m *Model) goForwardInView() {
	switch m.activeView {
	case taskListView:
		m.activeView = taskLogView
	case taskLogView:
		m.activeView = inactiveTaskListView
	case inactiveTaskListView:
		m.activeView = taskListView
	case saveActiveTLView, manualTasklogEntryView:
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
}

func (m *Model) goBackwardInView() {
	switch m.activeView {
	case taskLogView:
		m.activeView = taskListView
	case taskListView:
		m.activeView = inactiveTaskListView
	case inactiveTaskListView:
		m.activeView = taskLogView
	case saveActiveTLView, manualTasklogEntryView:
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

func (m *Model) shiftTime(direction types.TimeShiftDirection, duration types.TimeShiftDuration) error {
	if m.activeView == editActiveTLView || m.activeView == saveActiveTLView || m.activeView == manualTasklogEntryView {
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

func (m *Model) handleRequestToGoBackOrQuit() bool {
	var shouldQuit bool
	switch m.activeView {
	case taskListView:
		fs := m.activeTasksList.FilterState()
		if fs == list.Filtering || fs == list.FilterApplied {
			m.activeTasksList.ResetFilter()
		} else {
			shouldQuit = true
		}
	case taskLogView:
		fs := m.taskLogList.FilterState()
		if fs == list.Filtering || fs == list.FilterApplied {
			m.taskLogList.ResetFilter()
		} else {
			m.activeView = taskListView
		}
	case inactiveTaskListView:
		fs := m.inactiveTasksList.FilterState()
		if fs == list.Filtering || fs == list.FilterApplied {
			m.inactiveTasksList.ResetFilter()
		} else {
			m.activeView = taskListView
		}
	case helpView:
		m.activeView = taskListView
	default:
		shouldQuit = true
	}

	return shouldQuit
}

func (m *Model) getCmdToReloadData() tea.Cmd {
	var cmd tea.Cmd
	switch m.activeView {
	case taskListView:
		cmd = fetchTasks(m.db, true)
	case taskLogView:
		cmd = fetchTLS(m.db)
		m.taskLogList.ResetSelected()
	case inactiveTaskListView:
		cmd = fetchTasks(m.db, false)
		m.inactiveTasksList.ResetSelected()
	}

	return cmd
}

func (m *Model) goToActiveTask() {
	if m.activeView != taskListView {
		return
	}

	if !m.trackingActive {
		m.message = "Nothing is being tracked right now"
		return
	}

	if m.activeTasksList.IsFiltered() {
		m.activeTasksList.ResetFilter()
	}
	activeIndex, ok := m.activeTaskIndexMap[m.activeTaskID]
	if !ok {
		m.message = genericErrorMsg
		return
	}

	m.activeTasksList.Select(activeIndex)
}

func (m *Model) handleRequestToSaveActiveTL() {
	m.activeView = editActiveTLView
	m.trackingFocussedField = entryBeginTS
	m.trackingInputs[entryBeginTS].SetValue(m.activeTLBeginTS.Format(timeFormat))
	m.trackingInputs[m.trackingFocussedField].Focus()
}

func (m *Model) handleRequestToCreateManualTL() {
	m.activeView = manualTasklogEntryView
	m.tasklogSaveType = tasklogInsert
	m.trackingFocussedField = entryBeginTS
	currentTime := time.Now()
	currentTimeStr := currentTime.Format(timeFormat)

	m.trackingInputs[entryBeginTS].SetValue(currentTimeStr)
	m.trackingInputs[entryEndTS].SetValue(currentTimeStr)

	for i := range m.trackingInputs {
		m.trackingInputs[i].Blur()
	}
	m.trackingInputs[m.trackingFocussedField].Focus()
}

func (m *Model) getCmdToDeactivateTask() tea.Cmd {
	if m.activeTasksList.IsFiltered() {
		m.message = removeFilterMsg
		return nil
	}

	if m.trackingActive {
		m.message = "Cannot deactivate a task being tracked; stop tracking and try again."
		return nil
	}

	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = msgCouldntSelectATask
		return nil
	}

	return updateTaskActiveStatus(m.db, task, false)
}

func (m *Model) getCmdToDeleteTL() tea.Cmd {
	entry, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
	if !ok {
		m.message = "Couldn't delete task log entry"
		return nil
	}
	return deleteTL(m.db, &entry)
}

func (m *Model) getCmdToActivateDeactivatedTask() tea.Cmd {
	if m.inactiveTasksList.IsFiltered() {
		m.message = removeFilterMsg
		return nil
	}

	task, ok := m.inactiveTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = genericErrorMsg
		return nil
	}

	return updateTaskActiveStatus(m.db, task, true)
}

func (m *Model) getCmdToStartTracking() tea.Cmd {
	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = genericErrorMsg
		return nil
	}

	m.changesLocked = true
	m.activeTLBeginTS = time.Now()
	return toggleTracking(m.db, task.ID, m.activeTLBeginTS, m.activeTLEndTS, "")
}

func (m *Model) handleRequestToStopTracking() {
	m.activeView = saveActiveTLView
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

func (m *Model) handleRequestToCreateTask() {
	if m.activeTasksList.IsFiltered() {
		m.message = removeFilterMsg
		return
	}

	m.activeView = taskInputView
	m.taskInputFocussedField = summaryField
	m.taskInputs[summaryField].Focus()
	m.taskMgmtContext = taskCreateCxt
}

func (m *Model) handleRequestToUpdateTask() {
	if m.activeTasksList.IsFiltered() {
		m.message = removeFilterMsg
		return
	}

	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = genericErrorMsg
		return
	}

	m.activeView = taskInputView
	m.taskInputFocussedField = summaryField
	m.taskInputs[summaryField].Focus()
	m.taskInputs[summaryField].SetValue(task.Summary)
	m.taskMgmtContext = taskUpdateCxt
}

func (m *Model) handleRequestToScrollVPUp() {
	if m.helpVP.AtTop() {
		return
	}
	m.helpVP.LineUp(viewPortMoveLineCount)
}

func (m *Model) handleRequestToScrollVPDown() {
	if m.helpVP.AtBottom() {
		return
	}
	m.helpVP.LineDown(viewPortMoveLineCount)
}

func (m *Model) handleWindowResizing(msg tea.WindowSizeMsg) {
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
}

func (m *Model) handleTasksFetchedMsg(msg tasksFetchedMsg) tea.Cmd {
	if msg.err != nil {
		m.message = "error fetching tasks : " + msg.err.Error()
		return nil
	}

	var cmd tea.Cmd
	switch msg.active {
	case true:
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
		cmd = fetchActiveTask(m.db)

	case false:
		inactiveTasks := make([]list.Item, len(msg.tasks))
		for i, inactiveTask := range msg.tasks {
			inactiveTask.UpdateTitle()
			inactiveTask.UpdateDesc()
			inactiveTasks[i] = &inactiveTask
		}
		m.inactiveTasksList.SetItems(inactiveTasks)
	}

	return cmd
}

func (m *Model) handleManualTLInsertedMsg(msg manualTLInsertedMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = msg.err.Error()
		return nil
	}
	for i := range m.trackingInputs {
		m.trackingInputs[i].SetValue("")
	}
	task, ok := m.activeTaskMap[msg.taskID]

	var cmds []tea.Cmd
	if ok {
		cmds = append(cmds, updateTaskRep(m.db, task))
	}
	cmds = append(cmds, fetchTLS(m.db))

	return cmds
}

func (m *Model) handleTLSFetchedMsg(msg tLsFetchedMsg) {
	if msg.err != nil {
		m.message = msg.err.Error()
		return
	}

	items := make([]list.Item, len(msg.entries))
	for i, e := range msg.entries {
		e.UpdateTitle()
		e.UpdateDesc()
		items[i] = e
	}
	m.taskLogList.SetItems(items)
}

func (m *Model) handleActiveTaskFetchedMsg(msg activeTaskFetchedMsg) {
	if msg.err != nil {
		m.message = msg.err.Error()
		return
	}

	if msg.noneActive {
		m.lastTrackingChange = trackingFinished
		return
	}

	m.activeTaskID = msg.activeTaskID
	m.lastTrackingChange = trackingStarted
	m.activeTLBeginTS = msg.beginTs
	activeTask, ok := m.activeTaskMap[m.activeTaskID]
	if ok {
		activeTask.TrackingActive = true
		activeTask.UpdateTitle()

		// go to tracked item on startup
		activeIndex, aOk := m.activeTaskIndexMap[msg.activeTaskID]
		if aOk {
			m.activeTasksList.Select(activeIndex)
		}
	}
	m.trackingActive = true
}

func (m *Model) handleTrackingToggledMsg(msg trackingToggledMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = msg.err.Error()
		m.trackingActive = false
		return nil
	}

	m.changesLocked = false

	task, ok := m.activeTaskMap[msg.taskID]

	if !ok {
		m.message = genericErrorMsg
		return nil
	}

	var cmds []tea.Cmd
	switch msg.finished {
	case true:
		m.lastTrackingChange = trackingFinished
		task.TrackingActive = false
		m.trackingActive = false
		m.activeTaskID = -1
		cmds = append(cmds, updateTaskRep(m.db, task))
		cmds = append(cmds, fetchTLS(m.db))
	case false:
		m.lastTrackingChange = trackingStarted
		task.TrackingActive = true
		m.trackingActive = true
		m.activeTaskID = msg.taskID
	}

	task.UpdateTitle()

	return cmds
}

func (m *Model) handleTLDeleted(msg tLDeletedMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = "error deleting entry: " + msg.err.Error()
		return nil
	}

	var cmds []tea.Cmd
	task, ok := m.activeTaskMap[msg.entry.TaskID]
	if ok {
		cmds = append(cmds, updateTaskRep(m.db, task))
	}
	cmds = append(cmds, fetchTLS(m.db))

	return cmds
}

func (m *Model) handleActiveTLDeletedMsg(msg activeTaskLogDeletedMsg) {
	if msg.err != nil {
		m.message = fmt.Sprintf("Error deleting active log entry: %s", msg.err)
		return
	}

	activeTask, ok := m.activeTaskMap[m.activeTaskID]
	if !ok {
		m.message = genericErrorMsg
		return
	}

	activeTask.TrackingActive = false
	activeTask.UpdateTitle()
	m.lastTrackingChange = trackingFinished
	m.trackingActive = false
	m.activeTaskID = -1
}
