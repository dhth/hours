package theme

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/valid-with-entire-config.json
var validThemeWithEntireConfig []byte

//go:embed testdata/valid-with-partial-config.json
var validThemeWithPartialConfig []byte

//go:embed testdata/invalid-with-entire-config.json
var invalidThemeWithEntireConfig []byte

//go:embed testdata/malformed-json.json
var invalidThemeMalformedJSON []byte

//go:embed testdata/invalid-schema.json
var invalidThemeInvalidSchema []byte

//go:embed testdata/invalid-data.json
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
			defaultTheme := Default()
			err := json.Unmarshal(tt.themeBytes, &defaultTheme)
			require.NoError(t, err)
			// WHEN
			invalidColors := getInvalidColors(defaultTheme)

			// THEN
			assert.Len(t, invalidColors, tt.expectedNumInvalid)
		})
	}
}

func TestLoadCustomLoadsFullThemeCorrectly(t *testing.T) {
	// GIVEN
	// WHEN
	customTheme, err := loadCustom(validThemeWithEntireConfig)

	// THEN
	require.NoError(t, err)
	snaps.MatchStandaloneYAML(t, customTheme)
}

func TestLoadCustomPartialThemeCorrectly(t *testing.T) {
	// GIVEN
	// WHEN
	customTheme, err := loadCustom(validThemeWithPartialConfig)

	// THEN
	require.NoError(t, err)
	snaps.MatchStandaloneYAML(t, customTheme)
}

func TestLoadCustomHandlesFailuresCorrectly(t *testing.T) {
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
			_, err := loadCustom(tt.input)

			// THEN
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
