package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/hours/internal/types"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

var referenceTime = time.Date(2025, 8, 16, 9, 0, 0, 0, time.UTC)

func TestGetDurationValidityContext(t *testing.T) {
	testCases := []struct {
		name             string
		beginTS          string
		endTS            string
		expectedCtx      string
		expectedValidity tlFormValidity
	}{
		// success cases
		{
			name:             "less than an hour",
			beginTS:          "2025/08/08 00:40",
			endTS:            "2025/08/08 00:48",
			expectedCtx:      "You're recording 8m",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "exact hour",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 01:00",
			expectedCtx:      "You're recording 1h",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "hours and minutes",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 02:30",
			expectedCtx:      "You're recording 2h 30m",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "across day boundary",
			beginTS:          "2025/08/08 23:30",
			endTS:            "2025/08/09 00:15",
			expectedCtx:      "You're recording 45m",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "exactly at 8h threshold",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 08:00",
			expectedCtx:      "You're recording 8h",
			expectedValidity: tlSubmitOk,
		},
		{
			name:             "> 8h threshold",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 08:01",
			expectedCtx:      "You're recording 8h 1m",
			expectedValidity: tlSubmitWarn,
		},
		{
			name:             "very long duration",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/09 02:00",
			expectedCtx:      "You're recording 26h",
			expectedValidity: tlSubmitWarn,
		},
		// failure cases
		{
			name:             "empty begin time",
			beginTS:          "",
			endTS:            "2025/08/08 00:10",
			expectedCtx:      "Begin time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "empty end time",
			beginTS:          "2025/08/08 00:10",
			endTS:            "",
			expectedCtx:      "End time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "begin time as whitespace only",
			beginTS:          "   ",
			endTS:            "2025/08/08 00:10",
			expectedCtx:      "Begin time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "end time as whitespace only",
			beginTS:          "2025/08/08 00:10",
			endTS:            "   ",
			expectedCtx:      "End time is empty",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "invalid begin ts",
			beginTS:          "2025-08-08 00:10",
			endTS:            "2025/08/08 00:20",
			expectedCtx:      "Begin time is invalid",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "invalid end format",
			beginTS:          "2025/08/08 00:10",
			endTS:            "08-08-2025 00:20",
			expectedCtx:      "End time is invalid",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "end before start",
			beginTS:          "2025/08/08 01:00",
			endTS:            "2025/08/08 00:59",
			expectedCtx:      "End time is before begin time",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "zero duration",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 00:00",
			expectedCtx:      "You're recording no time, change begin and/or end time",
			expectedValidity: tlSubmitErr,
		},
		{
			name:             "one minute duration",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 00:01",
			expectedCtx:      "You're recording 1m",
			expectedValidity: tlSubmitOk,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gotCtx, gotValidity := getDurationValidityContext(tt.beginTS, tt.endTS)

			assert.Equal(t, tt.expectedCtx, gotCtx)
			assert.Equal(t, tt.expectedValidity, gotValidity)
		})
	}
}

// TODO: the following tests rely a lot on the internal details of the model, which works okay for basic snapshot tests.
// But a refactoring would be needed for more comprehensive tests.
// https://pkg.go.dev/github.com/charmbracelet/x/exp/teatest could be an option for proper E2E tests
func TestTaskListViewEmpty(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewWithTasks(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	task1 := createTestTask(1, "Implement feature A", true, false, m.timeProvider)
	task2 := createTestTask(2, "Fix bug in module B", true, true, m.timeProvider)
	task3 := createTestTask(3, "Write documentation", true, false, m.timeProvider)

	m.taskMap[1] = task1
	m.taskMap[2] = task2
	m.taskMap[3] = task3

	items := []list.Item{task1, task2, task3}
	m.activeTasksList.SetItems(items)

	m.trackingActive = true
	m.activeTaskID = 2
	m.activeTLBeginTS = referenceTime

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskLogViewEmpty(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskLogView

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskLogViewWithEntries(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskLogView

	entry1 := createTestTaskLogEntry(1, 1, "Implement feature A", m.timeProvider)
	entry2 := createTestTaskLogEntry(2, 2, "Fix bug in module B", m.timeProvider)
	entry3 := createTestTaskLogEntry(3, 1, "Implement feature A", m.timeProvider)

	items := []list.Item{entry1, entry2, entry3}
	m.taskLogList.SetItems(items)

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestEmptyInactiveTaskListView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = inactiveTaskListView

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestInactiveTaskListViewWithTasks(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = inactiveTaskListView

	task1 := createTestTask(4, "Archived feature", false, false, m.timeProvider)
	task2 := createTestTask(5, "Completed bug fix", false, false, m.timeProvider)

	items := []list.Item{task1, task2}
	m.inactiveTasksList.SetItems(items)

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestCreateTaskViewWithNoInput(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskInputView
	m.taskMgmtContext = taskCreateCxt

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestCreateTaskView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskInputView
	m.taskMgmtContext = taskCreateCxt
	m.taskInputs[summaryField].SetValue("a new task")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestUpdateTaskView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskInputView
	m.taskMgmtContext = taskUpdateCxt
	m.taskInputs[summaryField].SetValue("a task to be updated")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestFinishActiveTLView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = finishActiveTLView

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 09:00")
	m.tLInputs[entryEndTS].SetValue("2025/08/17 10:30")
	m.tLCommentInput.SetValue("Test comment for finishing task")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestEditActiveTLView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = editActiveTLView

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 09:00")
	m.tLCommentInput.SetValue("Updated comment")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestManualTasklogEntryView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = manualTasklogEntryView
	m.tasklogSaveType = tasklogInsert

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 09:00")
	m.tLInputs[entryEndTS].SetValue("2025/08/17 10:30")
	m.tLCommentInput.SetValue("Manual task log entry")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestEditSavedTLView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = editSavedTLView
	m.tasklogSaveType = tasklogUpdate

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 09:00")
	m.tLInputs[entryEndTS].SetValue("2025/08/17 10:30")
	m.tLCommentInput.SetValue("Edited saved task log")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestHelpView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = helpView
	m.helpVPReady = true

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestInsufficientDimensionsView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = insufficientDimensionsView
	m.terminalWidth = 50
	m.terminalHeight = 20

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestFinishActiveTLViewWhereEndTimeBeforeBeginTime(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = finishActiveTLView

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 10:30")
	m.tLInputs[entryEndTS].SetValue("2025/08/17 09:00")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestFinishActiveTLViewWhereNoTimeTracked(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = finishActiveTLView

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 10:30")
	m.tLInputs[entryEndTS].SetValue("2025/08/17 10:30")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestFinishActiveTLViewWithWarningContext(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = finishActiveTLView

	m.tLInputs[entryBeginTS].SetValue("2025/08/17 09:00")
	m.tLInputs[entryEndTS].SetValue("2025/08/17 18:30")

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewWithInfoContext(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true
	task := createTestTask(1, "Implement feature A", true, false, m.timeProvider)

	m.taskMap[1] = task

	items := []list.Item{task}
	m.activeTasksList.SetItems(items)
	m.message = userMsg{
		value:      "Task created successfully",
		kind:       userMsgInfo,
		framesLeft: 2,
	}

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewWithErrorMessage(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true
	m.message = userMsg{
		value:      "Error: Something went wrong",
		kind:       userMsgErr,
		framesLeft: 2,
	}

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewDebugMode(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.debug = true
	m.showHelpIndicator = false
	m.activeView = taskListView
	m.tasksFetched = true

	task1 := createTestTask(1, "Implement feature A", true, false, m.timeProvider)
	task2 := createTestTask(2, "Fix bug in module B", true, false, m.timeProvider)
	task3 := createTestTask(3, "Write documentation", true, false, m.timeProvider)

	m.taskMap[1] = task1
	m.taskMap[2] = task2
	m.taskMap[3] = task3

	items := []list.Item{task1, task2, task3}
	m.activeTasksList.SetItems(items)

	// WHEN
	result := m.View()

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func createTestModel() Model {
	theme := DefaultTheme()
	style := NewStyle(theme)

	testTimeProvider := types.TestTimeProvider{FixedTime: referenceTime}
	m := InitialModel(nil, style, testTimeProvider, false, logFramesConfig{})

	msg := tea.WindowSizeMsg{
		Width:  minWidthNeeded,
		Height: minHeightNeeded,
	}
	m.handleWindowResizing(msg)

	return m
}

func createTestTask(id int, summary string, active bool, trackingActive bool, tp types.TimeProvider) *types.Task {
	taskUpdateTime := referenceTime.Add(-3 * time.Hour)
	task := &types.Task{
		ID:             id,
		Summary:        summary,
		CreatedAt:      taskUpdateTime,
		UpdatedAt:      taskUpdateTime,
		TrackingActive: trackingActive,
		SecsSpent:      0,
		Active:         active,
	}

	task.UpdateListTitle()
	task.UpdateListDesc(tp)

	return task
}

func createTestTaskLogEntry(id int, taskID int, taskSummary string, tp types.TimeProvider) *types.TaskLogEntry {
	comment := "Test work on task"
	entryEndTime := referenceTime.Add(-1 * time.Hour)

	entry := &types.TaskLogEntry{
		ID:          id,
		TaskID:      taskID,
		TaskSummary: taskSummary,
		BeginTS:     entryEndTime.Add(-90 * time.Minute),
		EndTS:       entryEndTime,
		SecsSpent:   5400,
		Comment:     &comment,
	}

	entry.UpdateListTitle()
	entry.UpdateListDesc(tp)

	return entry
}
