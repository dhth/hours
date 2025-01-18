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
	beginTS, err := time.ParseInLocation(timeFormat, m.tLInputs[entryBeginTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}

	commentValue := strings.TrimSpace(m.tLCommentInput.Value())
	var comment *string
	if commentValue != "" {
		comment = &commentValue
	}

	m.activeView = taskListView
	return updateActiveTL(m.db, beginTS, comment)
}

func (m *Model) getCmdToFinishTrackingActiveTL() tea.Cmd {
	beginTS, err := time.ParseInLocation(timeFormat, m.tLInputs[entryBeginTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}
	m.activeTLBeginTS = beginTS

	endTS, err := time.ParseInLocation(timeFormat, m.tLInputs[entryEndTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}
	m.activeTLEndTS = endTS

	if m.activeTLEndTS.Sub(m.activeTLBeginTS).Seconds() <= 0 {
		m.message = "time spent needs to be positive"
		return nil
	}

	commentValue := strings.TrimSpace(m.tLCommentInput.Value())
	var comment *string
	if commentValue != "" {
		comment = &commentValue
	}

	for i := range m.tLInputs {
		m.tLInputs[i].SetValue("")
	}
	m.tLCommentInput.SetValue("")
	m.activeView = taskListView

	return toggleTracking(m.db, m.activeTaskID, m.activeTLBeginTS, m.activeTLEndTS, comment)
}

func (m *Model) getCmdToSaveOrUpdateTL() tea.Cmd {
	beginTS, err := time.ParseInLocation(timeFormat, m.tLInputs[entryBeginTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}

	endTS, err := time.ParseInLocation(timeFormat, m.tLInputs[entryEndTS].Value(), time.Local)
	if err != nil {
		m.message = err.Error()
		return nil
	}

	if endTS.Sub(beginTS).Seconds() <= 0 {
		m.message = "time spent needs to be positive"
		return nil
	}

	commentValue := strings.TrimSpace(m.tLCommentInput.Value())
	var comment *string
	if commentValue != "" {
		comment = &commentValue
	}

	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = "Something went wrong"
		return nil
	}
	if m.tasklogSaveType != tasklogInsert {
		return nil
	}

	m.blurTLTrackingInputs()
	m.tLCommentInput.SetValue("")
	m.activeTLComment = nil
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
	case finishActiveTLView:
		m.activeView = taskListView
		m.tLCommentInput.SetValue("")
	case manualTasklogEntryView:
		if m.tasklogSaveType == tasklogInsert {
			m.activeView = taskListView
		}
		for i := range m.tLInputs {
			m.tLInputs[i].SetValue("")
			m.tLCommentInput.SetValue("")
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
	case editActiveTLView:
		switch m.trackingFocussedField {
		case entryBeginTS:
			m.trackingFocussedField = entryComment
			m.tLInputs[entryBeginTS].Blur()
			m.tLCommentInput.Focus()
		case entryComment:
			m.trackingFocussedField = entryBeginTS
			m.tLInputs[entryBeginTS].Focus()
			m.tLCommentInput.Blur()
		}
	case finishActiveTLView, manualTasklogEntryView:
		switch m.trackingFocussedField {
		case entryBeginTS:
			m.trackingFocussedField = entryEndTS
			m.tLInputs[entryBeginTS].Blur()
			m.tLInputs[entryEndTS].Focus()
		case entryEndTS:
			m.trackingFocussedField = entryComment
			m.tLInputs[entryEndTS].Blur()
			m.tLCommentInput.Focus()
		case entryComment:
			m.trackingFocussedField = entryBeginTS
			m.tLCommentInput.Blur()
			m.tLInputs[entryBeginTS].Focus()
		}
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
	case editActiveTLView:
		switch m.trackingFocussedField {
		case entryBeginTS:
			m.trackingFocussedField = entryComment
			m.tLCommentInput.Focus()
			m.tLInputs[entryBeginTS].Blur()
		case entryComment:
			m.trackingFocussedField = entryBeginTS
			m.tLInputs[entryBeginTS].Focus()
			m.tLCommentInput.Blur()
		}
	case finishActiveTLView, manualTasklogEntryView:
		switch m.trackingFocussedField {
		case entryBeginTS:
			m.trackingFocussedField = entryComment
			m.tLCommentInput.Focus()
			m.tLInputs[entryBeginTS].Blur()
		case entryEndTS:
			m.trackingFocussedField = entryBeginTS
			m.tLInputs[entryBeginTS].Focus()
			m.tLInputs[entryEndTS].Blur()
		case entryComment:
			m.trackingFocussedField = entryEndTS
			m.tLInputs[entryEndTS].Focus()
			m.tLCommentInput.Blur()
		}
	}
}

func (m *Model) shiftTime(direction types.TimeShiftDirection, duration types.TimeShiftDuration) error {
	switch m.trackingFocussedField {
	case entryBeginTS, entryEndTS:
		ts, err := time.ParseInLocation(timeFormat, m.tLInputs[m.trackingFocussedField].Value(), time.Local)
		if err != nil {
			return err
		}

		newTs := types.GetShiftedTime(ts, direction, duration)

		m.tLInputs[m.trackingFocussedField].SetValue(newTs.Format(timeFormat))
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
	case taskLogDetailsView:
		m.activeView = taskLogView
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
	activeIndex, ok := m.taskIndexMap[m.activeTaskID]
	if !ok {
		m.message = genericErrorMsg
		return
	}

	m.activeTasksList.Select(activeIndex)
}

func (m *Model) handleRequestToEditActiveTL() {
	m.activeView = editActiveTLView
	m.tLInputs[entryBeginTS].SetValue(m.activeTLBeginTS.Format(timeFormat))
	if m.activeTLComment != nil {
		m.tLCommentInput.SetValue(*m.activeTLComment)
	} else {
		m.tLCommentInput.SetValue("")
	}

	m.blurTLTrackingInputs()
	m.tLInputs[entryBeginTS].Focus()
	m.trackingFocussedField = entryBeginTS
}

func (m *Model) handleRequestToCreateManualTL() {
	m.activeView = manualTasklogEntryView
	m.tasklogSaveType = tasklogInsert
	currentTime := time.Now()
	currentTimeStr := currentTime.Format(timeFormat)

	m.tLInputs[entryBeginTS].SetValue(currentTimeStr)
	m.tLInputs[entryEndTS].SetValue(currentTimeStr)

	m.blurTLTrackingInputs()
	m.trackingFocussedField = entryBeginTS
	m.tLInputs[entryBeginTS].Focus()
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
	m.activeTLBeginTS = time.Now().Truncate(time.Second)
	return toggleTracking(m.db, task.ID, m.activeTLBeginTS, m.activeTLEndTS, nil)
}

func (m *Model) handleRequestToStopTracking() {
	m.activeView = finishActiveTLView
	m.activeTLEndTS = time.Now()

	beginTimeStr := m.activeTLBeginTS.Format(timeFormat)
	currentTimeStr := m.activeTLEndTS.Format(timeFormat)

	m.tLInputs[entryBeginTS].SetValue(beginTimeStr)
	m.tLInputs[entryEndTS].SetValue(currentTimeStr)
	if m.activeTLComment != nil {
		m.tLCommentInput.SetValue(*m.activeTLComment)
	}
	m.trackingFocussedField = entryComment

	m.blurTLTrackingInputs()
	m.tLCommentInput.Focus()
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
	switch m.activeView {
	case helpView:
		if m.helpVP.AtTop() {
			return
		}
		m.helpVP.LineUp(viewPortMoveLineCount)
	case taskLogDetailsView:
		if m.tLDetailsVP.AtTop() {
			return
		}
		m.tLDetailsVP.LineUp(viewPortMoveLineCount)
	default:
		return
	}
}

func (m *Model) handleRequestToScrollVPDown() {
	switch m.activeView {
	case helpView:
		if m.helpVP.AtBottom() {
			return
		}
		m.helpVP.LineDown(viewPortMoveLineCount)
	case taskLogDetailsView:
		if m.tLDetailsVP.AtBottom() {
			return
		}
		m.tLDetailsVP.LineDown(viewPortMoveLineCount)
	default:
		return
	}
}

func (m *Model) handleRequestToViewTLDetails() {
	if len(m.taskLogList.Items()) == 0 {
		return
	}

	tl, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
	if !ok {
		m.message = genericErrorMsg
		return
	}

	var taskDetails string
	task, tOk := m.taskMap[tl.TaskID]
	if tOk {
		taskDetails = task.Summary
	}

	timeSpentStr := types.HumanizeDuration(tl.SecsSpent)

	details := fmt.Sprintf(`Task: %s

%s → %s (%s)

---

%s
`, taskDetails,
		tl.BeginTS.Format(timeFormat),
		tl.EndTS.Format(timeFormat),
		timeSpentStr,
		tl.GetComment())

	m.tLDetailsVP.SetContent(details)
	m.activeView = taskLogDetailsView
}

func (m *Model) handleWindowResizing(msg tea.WindowSizeMsg) {
	w, h := listStyle.GetFrameSize()

	m.terminalWidth = msg.Width
	m.terminalHeight = msg.Height

	if msg.Width < minWidthNeeded || msg.Height < minHeightNeeded {
		if m.activeView != insufficientDimensionsView {
			m.lastViewBeforeInsufficientDims = m.activeView
			m.activeView = insufficientDimensionsView
		}
		return
	}

	if m.activeView == insufficientDimensionsView {
		m.activeView = m.lastViewBeforeInsufficientDims
	}

	m.taskLogList.SetWidth(msg.Width - w)
	m.taskLogList.SetHeight(msg.Height - h - 2)

	m.activeTasksList.SetWidth(msg.Width - w)
	m.activeTasksList.SetHeight(msg.Height - h - 2)

	m.inactiveTasksList.SetWidth(msg.Width - w)
	m.inactiveTasksList.SetHeight(msg.Height - h - 2)

	if !m.helpVPReady {
		m.helpVP = viewport.New(msg.Width-4, m.terminalHeight-7)
		m.helpVP.SetContent(helpText)
		m.helpVP.KeyMap.Up.SetEnabled(false)
		m.helpVP.KeyMap.Down.SetEnabled(false)
		m.helpVPReady = true
	} else {
		m.helpVP.Height = m.terminalHeight - 7
		m.helpVP.Width = msg.Width - 4
	}

	if !m.tLDetailsVPReady {
		m.tLDetailsVP = viewport.New(msg.Width-4, m.terminalHeight-6)
		m.tLDetailsVP.KeyMap.Up.SetEnabled(false)
		m.tLDetailsVP.KeyMap.Down.SetEnabled(false)
		m.tLDetailsVPReady = true
	} else {
		m.tLDetailsVP.Height = m.terminalHeight - 6
		m.tLDetailsVP.Width = msg.Width - 4
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
		m.taskMap = make(map[int]*types.Task)
		m.taskIndexMap = make(map[int]int)
		tasks := make([]list.Item, len(msg.tasks))
		for i, task := range msg.tasks {
			task.UpdateListTitle()
			task.UpdateListDesc()
			tasks[i] = &task
			m.taskMap[task.ID] = &task
			m.taskIndexMap[task.ID] = i
		}
		m.activeTasksList.SetItems(tasks)
		m.activeTasksList.Title = "Tasks"
		m.tasksFetched = true
		cmd = fetchActiveTask(m.db)

	case false:
		inactiveTasks := make([]list.Item, len(msg.tasks))
		for i, inactiveTask := range msg.tasks {
			inactiveTask.UpdateListTitle()
			inactiveTask.UpdateListDesc()
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

	task, ok := m.taskMap[msg.taskID]

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
		e.UpdateListTitle()
		e.UpdateListDesc()
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

	m.lastTrackingChange = trackingStarted
	m.activeTaskID = msg.activeTask.TaskID
	m.activeTLBeginTS = msg.activeTask.CurrentLogBeginTS
	m.activeTLComment = msg.activeTask.CurrentLogComment

	activeTask, ok := m.taskMap[m.activeTaskID]
	if ok {
		activeTask.TrackingActive = true
		activeTask.UpdateListTitle()

		// go to tracked item on startup
		activeIndex, aOk := m.taskIndexMap[msg.activeTask.TaskID]
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

	task, ok := m.taskMap[msg.taskID]

	if !ok {
		m.message = genericErrorMsg
		return nil
	}

	var cmds []tea.Cmd
	switch msg.finished {
	case true:
		m.lastTrackingChange = trackingFinished
		task.TrackingActive = false
		m.activeTLComment = nil
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

	task.UpdateListTitle()

	return cmds
}

func (m *Model) handleTLDeleted(msg tLDeletedMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = "error deleting entry: " + msg.err.Error()
		return nil
	}

	var cmds []tea.Cmd
	task, ok := m.taskMap[msg.entry.TaskID]
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

	activeTask, ok := m.taskMap[m.activeTaskID]
	if !ok {
		m.message = genericErrorMsg
		return
	}

	activeTask.TrackingActive = false
	activeTask.UpdateListTitle()
	m.lastTrackingChange = trackingFinished
	m.trackingActive = false
	m.activeTLComment = nil
	m.activeTaskID = -1
}
