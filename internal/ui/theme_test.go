package ui

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThemeDefaults(t *testing.T) {
	defaultTheme := DefaultTheme()
	customThemeOrig := Theme{
		ActiveTaskList: "#ff0000",
		TaskEntry:      "#0000ff",
		HelpMsg:        "#00ff00",
	}
	customThemeJSON, err := json.Marshal(customThemeOrig)
	if err != nil {
		t.Errorf("Error marshalling theme: %s", err)
	}
	customTheme, err := LoadTheme(customThemeJSON)
	if err != nil {
		t.Errorf("Error loading theme: %s", err)
	}
	assert.Equal(t, customTheme.DefaultBackground, defaultTheme.DefaultBackground)
	assert.Equal(t, customTheme.ActiveTaskList, customThemeOrig.ActiveTaskList)
	assert.Equal(t, customTheme.InactiveTaskList, defaultTheme.InactiveTaskList)
	assert.Equal(t, customTheme.TaskEntry, customThemeOrig.TaskEntry)
	assert.Equal(t, customTheme.TaskLogEntry, defaultTheme.TaskLogEntry)
	assert.Equal(t, customTheme.TaskLogList, defaultTheme.TaskLogList)
	assert.Equal(t, customTheme.Tracking, defaultTheme.Tracking)
	assert.Equal(t, customTheme.ActiveTask, defaultTheme.ActiveTask)
	assert.Equal(t, customTheme.ActiveTaskBeginTime, defaultTheme.ActiveTaskBeginTime)
	assert.Equal(t, customTheme.FormFieldName, defaultTheme.FormFieldName)
	assert.Equal(t, customTheme.FormHelp, defaultTheme.FormHelp)
	assert.Equal(t, customTheme.FormContext, defaultTheme.FormContext)
	assert.Equal(t, customTheme.ToolName, defaultTheme.ToolName)
	assert.Equal(t, customTheme.RecordsHeader, defaultTheme.RecordsHeader)
	assert.Equal(t, customTheme.RecordsFooter, defaultTheme.RecordsFooter)
	assert.Equal(t, customTheme.RecordsBorder, defaultTheme.RecordsBorder)
	assert.Equal(t, customTheme.InitialHelpMsg, defaultTheme.InitialHelpMsg)
	assert.Equal(t, customTheme.RecordsDateRange, defaultTheme.RecordsDateRange)
	assert.Equal(t, customTheme.RecordsHelp, defaultTheme.RecordsHelp)
	assert.Equal(t, customTheme.HelpMsg, customThemeOrig.HelpMsg)
	assert.Equal(t, customTheme.HelpHeader, defaultTheme.HelpHeader)
	assert.Equal(t, customTheme.HelpSection, defaultTheme.HelpSection)
	assert.Equal(t, customTheme.FallbackTask, defaultTheme.FallbackTask)
	assert.Equal(t, customTheme.Warning, defaultTheme.Warning)
	assert.Equal(t, customTheme.Tasks, defaultTheme.Tasks)
}
