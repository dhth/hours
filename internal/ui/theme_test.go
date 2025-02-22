package ui

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed static/valid-theme-with-entire-config.json
var validThemeJSONWithEntireConfig []byte

//go:embed static/valid-theme-with-partial-config.json
var validThemeJSONWithPartialConfig []byte

//go:embed static/malformed-json-theme.json
var malformedJSONTheme []byte

//go:embed static/invalid-schema-theme.json
var invalidSchemaTheme []byte

//go:embed static/invalid-data-theme.json
var invalidDataTheme []byte

func TestGetInvalidColors(t *testing.T) {
	testCases := []struct {
		name               string
		themeBytes         []byte
		expectedNumInvalid int
	}{
		// success
		{
			name:               "valid json with all key-values provided",
			themeBytes:         validThemeJSONWithEntireConfig,
			expectedNumInvalid: 0,
		},
		// failures
		{
			name:               "invalid data",
			themeBytes:         invalidDataTheme,
			expectedNumInvalid: 5,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			theme := DefaultTheme()
			err := json.Unmarshal(tt.themeBytes, &theme)
			require.NoError(t, err)
			// WHEN
			invalidColors := getInvalidColors(theme)

			// THEN
			assert.Len(t, invalidColors, tt.expectedNumInvalid)
		})
	}
}

func TestLoadTheme(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
		err   error
	}{
		// success
		{
			name:  "valid json with all key-values provided",
			input: validThemeJSONWithEntireConfig,
		},
		{
			name:  "valid json with some key-values provided",
			input: validThemeJSONWithPartialConfig,
		},
		// failures
		{
			name:  "malformed json",
			input: malformedJSONTheme,
			err:   errThemeFileIsInvalidJSON,
		},
		{
			name:  "invalid schema",
			input: invalidSchemaTheme,
			err:   ErrThemeFileHasInvalidSchema,
		},
		{
			name:  "invalid data",
			input: invalidDataTheme,
			err:   ErrThemeColorsAreInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			// WHEN
			_, err := LoadTheme(tt.input)

			// THEN
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestLoadThemeFallsBacktoDefaults(t *testing.T) {
	// GIVEN
	defaultTheme := DefaultTheme()
	customThemeOrig := Theme{
		ActiveTasks: "#ff0000",
		TaskEntry:   "#0000ff",
		HelpMsg:     "#00ff00",
	}
	customThemeBytes, err := json.Marshal(customThemeOrig)
	require.NoError(t, err)

	// WHEN
	customTheme, err := LoadTheme(customThemeBytes)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, customTheme.TitleForeground, defaultTheme.TitleForeground)
	assert.Equal(t, customTheme.ActiveTasks, customThemeOrig.ActiveTasks)
	assert.Equal(t, customTheme.InactiveTasks, defaultTheme.InactiveTasks)
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
	assert.Equal(t, customTheme.HelpPrimary, defaultTheme.HelpPrimary)
	assert.Equal(t, customTheme.HelpSecondary, defaultTheme.HelpSecondary)
	assert.Equal(t, customTheme.Tasks, defaultTheme.Tasks)
}
