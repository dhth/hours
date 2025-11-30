package ui

import (
	"testing"

	"github.com/dhth/hours/internal/ui/theme"
	"github.com/stretchr/testify/assert"
)

func TestGetDynamicStyle(t *testing.T) {
	// GIVEN
	thm := theme.Default()
	style := NewStyle(thm)
	input := "abcdefghi"

	// WHEN
	gota := style.getDynamicStyle(input)
	gotb := style.getDynamicStyle(input)

	// THEN
	// assert same style returned for the same string
	assert.Equal(t, gota.GetForeground(), gotb.GetForeground())
}

func TestGetDynamicStyleHandlesEmptyTaskList(t *testing.T) {
	// GIVEN
	thm := theme.Default()
	thm.Tasks = []string{}
	style := NewStyle(thm)
	input := "abcdefghi"

	// WHEN
	got := style.getDynamicStyle(input)

	// THEN
	assert.NotNil(t, got)
}
