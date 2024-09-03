package persistence

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite" // sqlite driver
)

func TestMigrationsAreSetupCorrectly(t *testing.T) {
	// GIVEN
	// WHEN
	migrations := getMigrations()

	// THEN
	for i := 2; i <= latestDBVersion; i++ {
		m, ok := migrations[i]
		if !ok {
			assert.True(t, ok, "couldn't get migration %d", i)
		}
		if m == "" {
			assert.NotEmpty(t, ok, "migration %d is empty", i)
		}
	}
}

func TestMigrationsWork(t *testing.T) {
	// GIVEN
	var testDB *sql.DB
	var err error
	testDB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Couldn't open database: %s", err.Error())
	}

	err = InitDB(testDB)
	if err != nil {
		t.Fatalf("Couldn't initialize database: %s", err.Error())
	}

	// WHEN
	err = UpgradeDB(testDB, 1)

	// THEN
	assert.NoError(t, err)
}

func TestRunMigrationFailsWhenGivenBadMigration(t *testing.T) {
	// GIVEN
	var testDB *sql.DB
	var err error
	testDB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Couldn't open database: %s", err.Error())
	}

	err = InitDB(testDB)
	if err != nil {
		t.Fatalf("Couldn't initialize database: %s", err.Error())
	}

	// WHEN
	query := "BAD SQL CODE;"
	migrateErr := runMigration(testDB, query, 1)

	// THEN
	assert.Error(t, migrateErr)
}
