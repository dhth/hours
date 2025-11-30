package ui

import (
	"testing"

	"github.com/dhth/hours/internal/ui/theme"
	"github.com/stretchr/testify/assert"
)

func TestGetDynamicStyle(t *testing.T) {
	theme := theme.Default()
	style := NewStyle(theme)
	input := "abcdefghi"
	gota := style.getDynamicStyle(input)
	gotb := style.getDynamicStyle(input)
	// assert same style returned for the same string
	assert.Equal(t, gota.GetForeground(), gotb.GetForeground())
}
