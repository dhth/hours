package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDynamicStyle(t *testing.T) {
	input := "abcdefghi"
	gota := getDynamicStyle(input)
	gotb := getDynamicStyle(input)
	// assert same style returned for the same string
	assert.Equal(t, gota.GetForeground(), gotb.GetForeground())
}
