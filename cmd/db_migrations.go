package cmd

import (
	"database/sql"
	"time"
)

const (
	latestDBVersion = 1 // only upgrade this after adding a migration in getMigrations
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
	//     ADD COLUMN a_col INTEGER NOT NULL DEFAULT 1;
	// `

	return migrations
}

func fetchLatestDBVersion(db *sql.DB) (dbVersionInfo, error) {
	row := db.QueryRow(`
SELEC id, version, created_at
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

func upgradeDBIfNeeded(db *sql.DB) {
	latestVersionInDB, versionErr := fetchLatestDBVersion(db)
	if versionErr != nil {
		die(`Couldn't get hours' latest database version. This is a fatal error; let %s
know about this via %s.

Error: %s`,
			author,
			repoIssuesUrl,
			versionErr)
	}

	if latestVersionInDB.version > latestDBVersion {
		die(`Looks like you downgraded hours. You should either delete hours'
database file (you will lose data by doing that), or upgrade hours to
the latest version.`)
	}

	if latestVersionInDB.version < latestDBVersion {
		upgradeDB(db, latestVersionInDB.version)
	}
}

func upgradeDB(db *sql.DB, currentVersion int) {
	migrations := getMigrations()
	for i := currentVersion + 1; i <= latestDBVersion; i++ {
		migrateQuery := migrations[i]
		migrateErr := runMigration(db, migrateQuery, i)
		if migrateErr != nil {
			die(`Something went wrong migrating hours' database to version %d. This is not
supposed to happen. You can try running hours by passing it a custom database
file path (using --dbpath; this will create a new database) to see if that fixes
things. If that works, you can either delete the previous database, or keep
using this new database (both are not ideal).

If you can, let %s know about this error via
%s.
Sorry for breaking the upgrade step!

---

Error: %s
`, i, author, repoIssuesUrl, migrateErr)
		}
	}
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
