package ui

import (
	"testing"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/dhth/hours/internal/domain"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/ui/theme"
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
			name:             "> 8h threshold",
			beginTS:          "2025/08/08 00:00",
			endTS:            "2025/08/08 08:01",
			expectedCtx:      "You're recording 8h 1m",
			expectedValidity: tlSubmitWarn,
		},
		// failure cases
		{
			name:             "empty begin time",
			beginTS:          "",
			endTS:            "2025/08/08 00:10",
			expectedCtx:      "Error: begin time is empty",
			expectedValidity: tlSubmitErr,
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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewWhenTrackingTaskWithTimeSpent(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	task := createTestTask(1, "Fix bug in module B", true, false, m.timeProvider)
	task.SecsSpent = 130 * 60
	m.taskMap[1] = task
	m.activeTasksList.SetItems([]list.Item{task})

	m.handleActiveTaskFetchedMsg(activeTaskFetchedMsg{
		activeTask: domain.ActiveTaskDetails{TaskID: 1, CurrentLogBeginTS: referenceTime.Add(-38 * time.Minute)},
	})

	// WHEN
	updated, _ := m.Update(timeTrackedTodayTickMsg{})
	result := stripANSI(updated.(Model).View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewWhenTrackingTaskWithNoTimeSpent(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	task := createTestTask(1, "Fix bug in module B", true, false, m.timeProvider)
	m.taskMap[1] = task
	m.activeTasksList.SetItems([]list.Item{task})

	m.handleActiveTaskFetchedMsg(activeTaskFetchedMsg{
		activeTask: domain.ActiveTaskDetails{TaskID: 1, CurrentLogBeginTS: referenceTime.Add(-38 * time.Minute)},
	})

	// WHEN
	updated, _ := m.Update(timeTrackedTodayTickMsg{})
	result := stripANSI(updated.(Model).View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewClearsTrackedTimeWhenTrackingIsFinished(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	task := createTestTask(1, "Fix bug in module B", true, false, m.timeProvider)
	task.SecsSpent = 130 * 60
	m.taskMap[1] = task
	m.activeTasksList.SetItems([]list.Item{task})

	m.handleActiveTaskFetchedMsg(activeTaskFetchedMsg{
		activeTask: domain.ActiveTaskDetails{TaskID: 1, CurrentLogBeginTS: referenceTime.Add(-38 * time.Minute)},
	})
	tracking, _ := m.Update(timeTrackedTodayTickMsg{})
	snaps.MatchStandaloneSnapshot(t, stripANSI(tracking.(Model).View().Content))

	// WHEN
	// The saved total stays at 2h 10m here. In the app, finishing also returns a
	// command that refetches the total from the DB; this test doesn't run commands.
	updated, _ := tracking.Update(trackingToggledMsg{taskID: 1, finished: true})
	result := stripANSI(updated.(Model).View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewClearsTrackedTimeWhenTrackingIsCancelled(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	task := createTestTask(1, "Fix bug in module B", true, false, m.timeProvider)
	task.SecsSpent = 130 * 60
	m.taskMap[1] = task
	m.activeTasksList.SetItems([]list.Item{task})

	m.handleActiveTaskFetchedMsg(activeTaskFetchedMsg{
		activeTask: domain.ActiveTaskDetails{TaskID: 1, CurrentLogBeginTS: referenceTime.Add(-38 * time.Minute)},
	})
	tracking, _ := m.Update(timeTrackedTodayTickMsg{})
	snaps.MatchStandaloneSnapshot(t, stripANSI(tracking.(Model).View().Content))

	// WHEN
	updated, _ := tracking.Update(activeTaskLogDeletedMsg{})
	result := stripANSI(updated.(Model).View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskListViewMovesTrackedTimeWhenTrackingIsSwitched(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskListView
	m.tasksFetched = true

	task1 := createTestTask(1, "Fix bug in module B", true, false, m.timeProvider)
	task1.SecsSpent = 130 * 60
	task2 := createTestTask(2, "Write documentation", true, false, m.timeProvider)
	m.taskMap[1] = task1
	m.taskMap[2] = task2
	m.activeTasksList.SetItems([]list.Item{task1, task2})

	m.handleActiveTaskFetchedMsg(activeTaskFetchedMsg{
		activeTask: domain.ActiveTaskDetails{TaskID: 1, CurrentLogBeginTS: referenceTime.Add(-38 * time.Minute)},
	})
	tracking, _ := m.Update(timeTrackedTodayTickMsg{})
	snaps.MatchStandaloneSnapshot(t, stripANSI(tracking.(Model).View().Content))

	// WHEN
	// The old task's saved total stays at 2h 10m here. In the app, switching also returns a
	// command that refetches the total from the DB; this test doesn't run commands.
	updated, _ := tracking.Update(activeTLSwitchedMsg{lastActiveTaskID: 1, currentlyActiveTaskID: 2, ts: referenceTime})
	result := stripANSI(updated.(Model).View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTaskLogViewEmpty(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskLogView

	// WHEN
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestEmptyInactiveTaskListView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = inactiveTaskListView

	// WHEN
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTimeTrackedTodayInFooter(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskLogView
	m.secsTrackedToday = 90 * 60

	entry1 := createTestTaskLogEntry(1, 1, "Implement feature A", m.timeProvider)
	m.taskLogList.SetItems([]list.Item{entry1})

	// WHEN
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestCreateTaskViewWithNoInput(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskInputView
	m.taskMgmtContext = taskCreateCxt

	// WHEN
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestFinishActiveTLViewHidesTrackingTimesInFooter(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.tasksFetched = true

	task := createTestTask(1, "Fix bug in module B", true, false, m.timeProvider)
	m.taskMap[1] = task
	m.activeTasksList.SetItems([]list.Item{task})

	m.trackingActive = true
	m.activeTaskID = 1
	m.activeTLBeginTS = referenceTime.Add(-38 * time.Minute)

	m.activeView = finishActiveTLView
	m.tLInputs[entryBeginTS].SetValue(m.activeTLBeginTS.Format(timeFormat))
	m.tLInputs[entryEndTS].SetValue(referenceTime.Format(timeFormat))

	// WHEN
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestHelpView(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = helpView
	m.helpVPReady = true

	// WHEN
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

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
	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func TestTerminalWidthResizingWorks(t *testing.T) {
	// GIVEN
	m := createTestModel()
	m.activeView = taskLogView

	entry1 := createTestTaskLogEntry(1, 1, "Implement feature A", m.timeProvider)
	entry2 := createTestTaskLogEntry(2, 2, "Fix bug in module B", m.timeProvider)
	entry3 := createTestTaskLogEntry(3, 1, "Implement feature A", m.timeProvider)

	items := []list.Item{entry1, entry2, entry3}
	m.taskLogList.SetItems(items)

	// WHEN
	msg := tea.WindowSizeMsg{
		Width:  minWidthNeeded,
		Height: minHeightNeeded,
	}
	m.handleWindowResizing(msg)

	result := stripANSI(m.View().Content)

	// THEN
	snaps.MatchStandaloneSnapshot(t, result)
}

func createTestModel() Model {
	defaultTheme := theme.Default()
	style := NewStyle(defaultTheme)

	testTimeProvider := types.TestTimeProvider{FixedTime: referenceTime}
	m := InitialModel(nil, style, testTimeProvider, false, logFramesConfig{})

	msg := tea.WindowSizeMsg{
		Width:  96,
		Height: minHeightNeeded,
	}
	m.handleWindowResizing(msg)

	return m
}

func createTestTask(id int, summary string, active bool, trackingActive bool, tp types.TimeProvider) *taskListItem {
	taskUpdateTime := referenceTime.Add(-3 * time.Hour)
	task := &taskListItem{
		Task: domain.Task{
			ID:        id,
			Summary:   summary,
			CreatedAt: taskUpdateTime,
			UpdatedAt: taskUpdateTime,
			SecsSpent: 0,
			Active:    active,
		},
		trackingActive: trackingActive,
	}

	task.updateListTitle()
	task.updateListDesc(nil, tp.Now())

	return task
}

func createTestTaskLogEntry(id int, taskID int, taskSummary string, tp types.TimeProvider) taskLogListItem {
	comment := "Test work on task"
	entryEndTime := referenceTime.Add(-1 * time.Hour)

	entry := taskLogListItem{
		TaskLogEntry: domain.TaskLogEntry{
			ID:          id,
			TaskID:      taskID,
			TaskSummary: taskSummary,
			BeginTS:     entryEndTime.Add(-90 * time.Minute),
			EndTS:       entryEndTime,
			SecsSpent:   5400,
			Comment:     &comment,
		},
	}

	entry.updateListTitle()
	entry.updateListDesc(tp)

	return entry
}
