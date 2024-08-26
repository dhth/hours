package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	latestDBVersion = 1 // only upgrade this after adding a migration in getMigrations
	timeFormat      = "2006/01/02 15:04"
)

var (
	ErrDBDowngraded          = errors.New("database downgraded")
	ErrDBMigrationFailed     = errors.New("database migration failed")
	ErrCouldntFetchDBVersion = errors.New("couldn't fetch version")
)

type dbVersionInfo struct {
	id        int
	version   int
	createdAt time.Time
}

func getMigrations() map[int]string {
	migrations := make(map[int]string)
	// these migrations should not be modified once released.
	// that is, migrations is an append-only map.

	// migrations[2] = `
	// ALTER TABLE task
	// ADD COLUMN new_col TEXT;
	// `

	return migrations
}

func fetchLatestDBVersion(db *sql.DB) (dbVersionInfo, error) {
	row := db.QueryRow(`
SELECT id, version, created_at
FROM db_versions
ORDER BY created_at DESC
LIMIT 1;
`)

	var dbVersion dbVersionInfo
	err := row.Scan(
		&dbVersion.id,
		&dbVersion.version,
		&dbVersion.createdAt,
	)

	return dbVersion, err
}

func UpgradeDBIfNeeded(db *sql.DB) error {
	latestVersionInDB, err := fetchLatestDBVersion(db)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCouldntFetchDBVersion, err.Error())
	}

	if latestVersionInDB.version > latestDBVersion {
		return fmt.Errorf("%w; debug info: version=%d, created at=%q)",
			ErrDBDowngraded,
			latestVersionInDB.version,
			latestVersionInDB.createdAt.Format(timeFormat),
		)
	}

	if latestVersionInDB.version < latestDBVersion {
		err = UpgradeDB(db, latestVersionInDB.version)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpgradeDB(db *sql.DB, currentVersion int) error {
	migrations := getMigrations()
	for i := currentVersion + 1; i <= latestDBVersion; i++ {
		migrateQuery := migrations[i]
		migrateErr := runMigration(db, migrateQuery, i)
		if migrateErr != nil {
			return fmt.Errorf("%w (version %d): %v", ErrDBMigrationFailed, i, migrateErr.Error())
		}
	}
	return nil
}

func runMigration(db *sql.DB, migrateQuery string, version int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(migrateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	tStmt, err := tx.Prepare(`
INSERT INTO db_versions (version, created_at)
VALUES (?, ?);
`)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	_, err = tStmt.Exec(version, time.Now().UTC())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
