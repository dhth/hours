package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationsAreSetupCorrectly(t *testing.T) {
	migrations := getMigrations()
	for i := 2; i <= latestDBVersion; i++ {
		m, ok := migrations[i]
		assert.True(t, ok)
		assert.NotEmpty(t, m)
	}
}
