package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/stretchr/testify/assert"
)

func TestHandleCopyTaskSummary(t *testing.T) {
	testCases := []struct {
		name            string
		setupModel      func() Model
		expectedMsg     string
		expectedMsgKind userMsgKind
	}{
		{
			name: "success - active task list",
			setupModel: func() Model {
				m := createTestModel()
				m.activeView = taskListView
				task := createTestTask(1, "Test task summary", true, false, m.timeProvider)
				m.taskMap[1] = task
				m.activeTasksList.SetItems([]list.Item{task})
				m.activeTasksList.Select(0)
				return m
			},
			expectedMsg:     "Copied to clipboard",
			expectedMsgKind: userMsgInfo,
		},
		{
			name: "success - inactive task list",
			setupModel: func() Model {
				m := createTestModel()
				m.activeView = inactiveTaskListView
				task := createTestTask(1, "Archived task", false, false, m.timeProvider)
				m.inactiveTasksList.SetItems([]list.Item{task})
				m.inactiveTasksList.Select(0)
				return m
			},
			expectedMsg:     "Copied to clipboard",
			expectedMsgKind: userMsgInfo,
		},
		{
			name: "no task selected - active task list",
			setupModel: func() Model {
				m := createTestModel()
				m.activeView = taskListView
				return m
			},
			expectedMsg:     "No task selected",
			expectedMsgKind: userMsgErr,
		},
		{
			name: "no task selected - inactive task list",
			setupModel: func() Model {
				m := createTestModel()
				m.activeView = inactiveTaskListView
				return m
			},
			expectedMsg:     "No task selected",
			expectedMsgKind: userMsgErr,
		},
		{
			name: "wrong view - task log view",
			setupModel: func() Model {
				m := createTestModel()
				m.activeView = taskLogView
				return m
			},
			expectedMsg:     "",
			expectedMsgKind: userMsgInfo,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupModel()
			m.handleCopyTaskSummary()

			assert.Equal(t, tt.expectedMsg, m.message.value)
			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsgKind, m.message.kind)
			}
		})
	}
}
