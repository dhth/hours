package persistence

import (
	"testing"
)

func TestMigrationsAreSetupCorrectly(t *testing.T) {
	migrations := getMigrations()
	for i := 2; i <= latestDBVersion; i++ {
		m, ok := migrations[i]
		if !ok {
			t.Errorf("couldn't get migration %d", i)
		}
		if m == "" {
			t.Errorf("migration %d is empty", i)
		}
	}
}
