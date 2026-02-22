package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aymanbagabas/go-osc52/v2"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	c "github.com/dhth/hours/internal/common"
	"github.com/dhth/hours/internal/types"
)

const (
	genericErrorMsg               = "Something went wrong"
	removeFilterMsg               = "Remove filter first"
	beginTsCannotBeInTheFutureMsg = "Begin timestamp cannot be in the future"
)

var suggestReloadingMsg = fmt.Sprintf("Something went wrong, please restart hours; let %s know about this error via %s.", c.Author, c.RepoIssuesURL)

func (m *Model) getCmdToCreateOrUpdateTask() tea.Cmd {
	if strings.TrimSpace(m.taskInputs[summaryField].Value()) == "" {
		m.message = errMsg("Task summary cannot be empty")
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
			m.message = errMsg("Something went wrong")
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
		m.message = errMsgQuick(err.Error())
		return nil
	}

	if beginTS.After(m.timeProvider.Now()) {
		m.message = errMsgQuick(beginTsCannotBeInTheFutureMsg)
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
	beginTS, endTS, err := types.ParseTaskLogTimes(m.tLInputs[entryBeginTS].Value(), m.tLInputs[entryEndTS].Value())
	if err != nil {
		return nil
	}

	m.activeTLBeginTS = beginTS
	m.activeTLEndTS = endTS

	commentValue := strings.TrimSpace(m.tLCommentInput.Value())
	var comment *string
	if commentValue != "" {
		comment = &commentValue
	}

	m.activeView = taskListView

	return toggleTracking(m.db, m.activeTaskID, m.activeTLBeginTS, m.activeTLEndTS, comment)
}

func (m *Model) getCmdToFinishActiveTLWithoutComment() tea.Cmd {
	now := m.timeProvider.Now().Truncate(time.Second)
	err := types.IsTaskLogDurationValid(m.activeTLBeginTS, now)

	if errors.Is(err, types.ErrDurationNotLongEnough) {
		m.message = infoMsg("Task log duration is too short to save; press <ctrl+x> if you want to discard it")
		return nil
	}

	if err != nil {
		m.message = errMsg(fmt.Sprintf("Error: %s", err.Error()))
		return nil
	}

	m.activeTLEndTS = now

	return toggleTracking(m.db, m.activeTaskID, m.activeTLBeginTS, m.activeTLEndTS, nil)
}

func (m *Model) getCmdToCreateOrEditTL() tea.Cmd {
	beginTS, endTS, err := types.ParseTaskLogTimes(m.tLInputs[entryBeginTS].Value(), m.tLInputs[entryEndTS].Value())
	if err != nil {
		return nil
	}

	commentValue := strings.TrimSpace(m.tLCommentInput.Value())
	var comment *string
	if commentValue != "" {
		comment = &commentValue
	}

	m.blurTLTrackingInputs()
	m.tLCommentInput.SetValue("")
	m.activeTLComment = nil

	var cmd tea.Cmd
	switch m.tasklogSaveType {
	case tasklogInsert:
		m.activeView = taskListView
		task, ok := m.activeTasksList.SelectedItem().(*types.Task)
		if !ok {
			m.message = errMsg(genericErrorMsg)
			return nil
		}
		cmd = insertManualTL(m.db, task.ID, beginTS, endTS, comment)
	case tasklogUpdate:
		m.activeView = taskLogView
		tl, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
		if !ok {
			m.message = errMsg(genericErrorMsg)
			return nil
		}
		cmd = editSavedTL(m.db, tl.ID, tl.TaskID, beginTS, endTS, comment)
	}

	return cmd
}

func (m *Model) handleEscapeInForms() {
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
	case editSavedTLView:
		m.activeView = taskLogView
	case moveTaskLogView:
		m.activeView = taskLogView
		m.targetTasksList.ResetFilter()
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
	case finishActiveTLView, manualTasklogEntryView, editSavedTLView:
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
	case finishActiveTLView, manualTasklogEntryView, editSavedTLView:
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
			m.activeView = taskLogView
		}
	case helpView:
		m.activeView = m.lastView
	case moveTaskLogView:
		m.activeView = taskLogView
		m.targetTasksList.ResetFilter()
	}

	return shouldQuit
}

func (m *Model) getCmdToReloadData() tea.Cmd {
	var cmd tea.Cmd
	switch m.activeView {
	case taskListView:
		cmd = fetchTasks(m.db, true)
	case taskLogView:
		cmd = fetchTLS(m.db, nil)
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
		m.message = errMsg("Nothing is being tracked right now")
		return
	}

	if m.activeTasksList.IsFiltered() {
		m.activeTasksList.ResetFilter()
	}
	activeIndex, ok := m.taskIndexMap[m.activeTaskID]
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return
	}

	m.activeTasksList.Select(activeIndex)
}

func (m *Model) handleRequestToEditActiveTL() {
	m.clearAllTaskLogInputs()
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
	m.clearAllTaskLogInputs()
	m.activeView = manualTasklogEntryView
	m.tasklogSaveType = tasklogInsert
	currentTime := m.timeProvider.Now()
	currentTimeStr := currentTime.Format(timeFormat)

	m.tLInputs[entryBeginTS].SetValue(currentTimeStr)
	m.tLInputs[entryEndTS].SetValue(currentTimeStr)

	m.blurTLTrackingInputs()
	m.trackingFocussedField = entryBeginTS
	m.tLInputs[entryBeginTS].Focus()
}

func (m *Model) handleRequestToEditSavedTL() {
	if len(m.taskLogList.Items()) == 0 {
		return
	}

	tl, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return
	}

	m.activeView = editSavedTLView
	m.tasklogSaveType = tasklogUpdate

	beginTimeStr := tl.BeginTS.Format(timeFormat)
	endTimeStr := tl.EndTS.Format(timeFormat)

	var comment string
	if tl.Comment != nil {
		comment = *tl.Comment
	}

	m.tLInputs[entryBeginTS].SetValue(beginTimeStr)
	m.tLInputs[entryEndTS].SetValue(endTimeStr)
	m.tLCommentInput.SetValue(comment)

	m.blurTLTrackingInputs()
	m.trackingFocussedField = entryBeginTS
	m.tLInputs[entryBeginTS].Focus()
}

func (m *Model) getCmdToDeactivateTask() tea.Cmd {
	if m.activeTasksList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return nil
	}

	if m.trackingActive {
		m.message = errMsg("Cannot deactivate a task being tracked; stop tracking and try again.")
		return nil
	}

	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(msgCouldntSelectATask)
		return nil
	}

	return updateTaskActiveStatus(m.db, task, false)
}

func (m *Model) getCmdToDeleteTL() tea.Cmd {
	entry, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
	if !ok {
		m.message = errMsg("Couldn't delete task log entry")
		return nil
	}
	return deleteTL(m.db, &entry)
}

func (m *Model) handleRequestToMoveTaskLog() tea.Cmd {
	if m.taskLogList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return nil
	}

	entry, ok := m.taskLogList.SelectedItem().(types.TaskLogEntry)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return nil
	}

	// Store the log entry details
	m.moveTLID = entry.ID
	m.moveOldTaskID = entry.TaskID
	m.moveSecsSpent = entry.SecsSpent

	// Initialize target list with active tasks, excluding current parent
	targetItems := []list.Item{}
	for i := range m.activeTasksList.Items() {
		task, ok := m.activeTasksList.Items()[i].(*types.Task)
		if !ok {
			continue
		}
		// Exclude the current parent task
		if task.ID != entry.TaskID {
			targetItems = append(targetItems, task)
		}
	}
	m.targetTasksList.SetItems(targetItems)

	m.activeView = moveTaskLogView
	return nil
}

func (m *Model) handleTargetTaskSelection() tea.Cmd {
	task, ok := m.targetTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return nil
	}

	return moveTaskLog(m.db, m.moveTLID, m.moveOldTaskID, task.ID, m.moveSecsSpent)
}

func (m *Model) getCmdToActivateDeactivatedTask() tea.Cmd {
	if m.inactiveTasksList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return nil
	}

	task, ok := m.inactiveTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return nil
	}

	return updateTaskActiveStatus(m.db, task, true)
}

func (m *Model) getCmdToStartTracking() tea.Cmd {
	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return nil
	}

	m.changesLocked = true
	m.activeTLBeginTS = m.timeProvider.Now().Truncate(time.Second)
	return toggleTracking(m.db, task.ID, m.activeTLBeginTS, m.activeTLEndTS, nil)
}

func (m *Model) handleRequestToStopTracking() {
	m.clearAllTaskLogInputs()
	m.activeView = finishActiveTLView
	m.activeTLEndTS = m.timeProvider.Now()

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

func (m *Model) getCmdToQuickSwitchTracking() tea.Cmd {
	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return nil
	}

	if task.ID == m.activeTaskID {
		return nil
	}

	if !m.trackingActive {
		m.changesLocked = true
		m.activeTLBeginTS = m.timeProvider.Now().Truncate(time.Second)
		return toggleTracking(m.db,
			task.ID,
			m.activeTLBeginTS,
			m.activeTLEndTS,
			nil,
		)
	}

	return quickSwitchActiveIssue(m.db, task.ID, m.timeProvider.Now())
}

func (m *Model) handleRequestToCreateTask() {
	if m.activeTasksList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return
	}

	m.activeView = taskInputView
	m.taskInputFocussedField = summaryField
	m.taskInputs[summaryField].Focus()
	m.taskMgmtContext = taskCreateCxt
}

func (m *Model) handleRequestToUpdateTask() {
	if m.activeTasksList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return
	}

	task, ok := m.activeTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return
	}

	m.activeView = taskInputView
	m.taskInputFocussedField = summaryField
	m.taskInputs[summaryField].Focus()
	m.taskInputs[summaryField].SetValue(task.Summary)
	m.taskMgmtContext = taskUpdateCxt
}

func (m *Model) handleCopyTaskSummary() {
	var selectedTask *types.Task
	var ok bool

	switch m.activeView {
	case taskListView:
		selectedTask, ok = m.activeTasksList.SelectedItem().(*types.Task)
	case inactiveTaskListView:
		selectedTask, ok = m.inactiveTasksList.SelectedItem().(*types.Task)
	default:
		return
	}

	if !ok || selectedTask == nil {
		m.message = errMsg("No task selected")
		return
	}

	osc52.New(selectedTask.Summary).WriteTo(os.Stderr)
	m.message = infoMsg("Copied to clipboard")
}

func (m *Model) handleRequestToScrollVPUp() {
	switch m.activeView {
	case helpView:
		if m.helpVP.AtTop() {
			return
		}
		m.helpVP.ScrollUp(viewPortMoveLineCount)
	case taskLogDetailsView:
		if m.tLDetailsVP.AtTop() {
			return
		}
		m.tLDetailsVP.ScrollUp(viewPortMoveLineCount)
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
		m.helpVP.ScrollDown(viewPortMoveLineCount)
	case taskLogDetailsView:
		if m.tLDetailsVP.AtBottom() {
			return
		}
		m.tLDetailsVP.ScrollDown(viewPortMoveLineCount)
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
		m.message = errMsg(genericErrorMsg)
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
	w, h := m.style.list.GetFrameSize()

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

	m.targetTasksList.SetWidth(msg.Width - w)
	m.targetTasksList.SetHeight(msg.Height - h - 2)

	if !m.helpVPReady {
		m.helpVP = viewport.New(msg.Width-4, m.terminalHeight-7)
		m.helpVP.SetContent(getHelpText(m.style))
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
		m.message = errMsg("Error fetching tasks : " + msg.err.Error())
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
			task.UpdateListDesc(m.timeProvider)
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
			inactiveTask.UpdateListDesc(m.timeProvider)
			inactiveTasks[i] = &inactiveTask
		}
		m.inactiveTasksList.SetItems(inactiveTasks)
	}

	return cmd
}

func (m *Model) handleManualTLInsertedMsg(msg manualTLInsertedMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = errMsg(msg.err.Error())
		return nil
	}

	task, ok := m.taskMap[msg.taskID]

	var cmds []tea.Cmd
	if ok {
		cmds = append(cmds, updateTaskRep(m.db, task))
	}
	cmds = append(cmds, fetchTLS(m.db, nil))

	return cmds
}

func (m *Model) handleSavedTLEditedMsg(msg savedTLEditedMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = errMsg(msg.err.Error())
		return nil
	}

	task, ok := m.taskMap[msg.taskID]

	var cmds []tea.Cmd
	if ok {
		cmds = append(cmds, updateTaskRep(m.db, task))
	}
	cmds = append(cmds, fetchTLS(m.db, &msg.tlID))

	return cmds
}

func (m *Model) handleTLSFetchedMsg(msg tLsFetchedMsg) {
	if msg.err != nil {
		m.message = errMsg(msg.err.Error())
		return
	}

	items := make([]list.Item, len(msg.entries))
	var indexToFocusOn *int
	var indexToFocusOnFound bool
	for i, e := range msg.entries {
		e.UpdateListTitle()
		e.UpdateListDesc(m.timeProvider)
		items[i] = e
		if !indexToFocusOnFound && msg.tlIDToFocusOn != nil && e.ID == *msg.tlIDToFocusOn {
			indexToFocusOn = &i
			indexToFocusOnFound = true
		}
	}
	m.taskLogList.SetItems(items)

	if indexToFocusOn != nil {
		m.taskLogList.Select(*indexToFocusOn)
	} else {
		m.taskLogList.Select(0)
	}
}

func (m *Model) handleActiveTaskFetchedMsg(msg activeTaskFetchedMsg) {
	if msg.err != nil {
		m.message = errMsg(msg.err.Error())
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
		m.message = errMsg(msg.err.Error())
		m.trackingActive = false
		return nil
	}

	m.changesLocked = false

	task, ok := m.taskMap[msg.taskID]

	if !ok {
		m.message = errMsg(genericErrorMsg)
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
		cmds = append(cmds, fetchTLS(m.db, nil))
	case false:
		m.lastTrackingChange = trackingStarted
		task.TrackingActive = true
		m.trackingActive = true
		m.activeTaskID = msg.taskID
	}

	task.UpdateListTitle()

	return cmds
}

func (m *Model) handleActiveTLSwitchedMsg(msg activeTLSwitchedMsg) tea.Cmd {
	if msg.err != nil {
		m.message = errMsg(msg.err.Error())
		return nil
	}

	lastActiveTask, ok := m.taskMap[msg.lastActiveTaskID]

	if !ok {
		m.message = errMsg(suggestReloadingMsg)
		return nil
	}

	lastActiveTask.TrackingActive = false
	lastActiveTask.UpdateListTitle()

	currentlyActiveTask, ok := m.taskMap[msg.currentlyActiveTaskID]

	if !ok {
		m.message = errMsg(suggestReloadingMsg)
		return nil
	}
	currentlyActiveTask.TrackingActive = true
	currentlyActiveTask.UpdateListTitle()

	m.activeTLComment = nil
	m.activeTaskID = msg.currentlyActiveTaskID
	m.activeTLBeginTS = msg.ts

	return fetchTLS(m.db, nil)
}

func (m *Model) handleTLDeleted(msg tLDeletedMsg) []tea.Cmd {
	if msg.err != nil {
		m.message = errMsg("Error deleting entry: " + msg.err.Error())
		return nil
	}

	var cmds []tea.Cmd
	task, ok := m.taskMap[msg.entry.TaskID]
	if ok {
		cmds = append(cmds, updateTaskRep(m.db, task))
	}
	cmds = append(cmds, fetchTLS(m.db, nil))

	return cmds
}

func (m *Model) handleActiveTLDeletedMsg(msg activeTaskLogDeletedMsg) {
	if msg.err != nil {
		m.message = errMsg(fmt.Sprintf("Error deleting active log entry: %s", msg.err))
		return
	}

	activeTask, ok := m.taskMap[m.activeTaskID]
	if !ok {
		m.message = errMsg(genericErrorMsg)
		return
	}

	activeTask.TrackingActive = false
	activeTask.UpdateListTitle()
	m.lastTrackingChange = trackingFinished
	m.trackingActive = false
	m.activeTLComment = nil
	m.activeTaskID = -1
}

func (m *Model) clearAllTaskLogInputs() {
	for i := range m.tLInputs {
		m.tLInputs[i].SetValue("")
	}
	m.tLCommentInput.SetValue("")
}
