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

func fetchDBVersion(db *sql.DB) (dbVersionInfo, error) {
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

func runMigration(db *sql.DB, migrateQuery string, version int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

func exitIfDBNeedsUpgrade(db *sql.DB) {
	latestVersion, versionErr := fetchDBVersion(db)
	if versionErr != nil {
		die("Couldn't get hours' latest database version. This is a fatal error; let @dhth know about this via %s.\nError: %s",
			repoIssuesUrl,
			versionErr)
	}

	if latestVersion.version < latestDBVersion {
		die("hours' database needs an upgrade. Run \"hours db upgrade\" to do so.")
	}

	if latestVersion.version > latestDBVersion {
		die("Looks like you downgraded hours. You should either delete hours' database file, or upgrade hours to the latest version.")
	}
}

func upgradeDB(db *sql.DB, currentVersion int) {
	migrations := getMigrations()
	for i := currentVersion + 1; i <= latestDBVersion; i++ {
		migrateQuery := migrations[i]
		migrateErr := runMigration(db, migrateQuery, i)
		if migrateErr != nil {
			die(`
Something went wrong migrating the database to version %d. This is not supposed
to happen.

You can try again by running the same command again. But, if that doesn't work,
the only way to recover for now is to delete the database file (the default
location can be checked via "hours -h"), and create a new one using "hours db
init".

If you can let, @dhth know about this error via
https://github.com/dhth/hours/issues. Sorry for breaking the upgrade step!
`, i)
		}
	}
}

func getDB(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbpath)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, err
}

func initDB(db *sql.DB) error {
	// these init queries cannot be changed
	// once hours is released; only further migrations
	// can be added, which are run via hours db upgrade
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS db_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    summary TEXT NOT NULL,
    secs_spent INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER,
    begin_ts TIMESTAMP NOT NULL,
    end_ts TIMESTAMP,
    secs_spent INTEGER NOT NULL DEFAULT 0,
    comment VARCHAR(255),
    active BOOLEAN NOT NULL,
    FOREIGN KEY(task_id) REFERENCES task(id)
);

CREATE TRIGGER IF NOT EXISTS prevent_duplicate_active_insert
BEFORE INSERT ON task_log
BEGIN
    SELECT CASE
        WHEN EXISTS (SELECT 1 FROM task_log WHERE active = 1)
        THEN RAISE(ABORT, 'Only one row with active=1 is allowed')
    END;
END;

INSERT INTO db_versions (version, created_at)
VALUES (1, ?);
`, time.Now().UTC())

	return err
}
