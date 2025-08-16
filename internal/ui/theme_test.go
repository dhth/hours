package ui

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed static/valid-with-entire-config.json
var validThemeWithEntireConfig []byte

//go:embed static/valid-with-partial-config.json
var validThemeWithPartialConfig []byte

//go:embed static/invalid-with-entire-config.json
var invalidThemeWithEntireConfig []byte

//go:embed static/malformed-json.json
var invalidThemeMalformedJSON []byte

//go:embed static/invalid-schema.json
var invalidThemeInvalidSchema []byte

//go:embed static/invalid-data.json
var invalidThemeInvalidData []byte

func TestGetInvalidColors(t *testing.T) {
	testCases := []struct {
		name               string
		themeBytes         []byte
		expectedNumInvalid int
	}{
		// success
		{
			name:               "valid json with all key-values provided",
			themeBytes:         validThemeWithEntireConfig,
			expectedNumInvalid: 0,
		},
		// failures
		{
			name:               "invalid data",
			themeBytes:         invalidThemeInvalidData,
			expectedNumInvalid: 5,
		},
		{
			name:               "invalid data with entire config",
			themeBytes:         invalidThemeWithEntireConfig,
			expectedNumInvalid: 42,
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
			input: validThemeWithEntireConfig,
		},
		{
			name:  "valid json with some key-values provided",
			input: validThemeWithPartialConfig,
		},
		// failures
		{
			name:  "malformed json",
			input: invalidThemeMalformedJSON,
			err:   errThemeFileIsInvalidJSON,
		},
		{
			name:  "invalid schema",
			input: invalidThemeInvalidSchema,
			err:   ErrThemeFileHasInvalidSchema,
		},
		{
			name:  "invalid data",
			input: invalidThemeInvalidData,
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
