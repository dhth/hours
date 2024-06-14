package cmd

import (
	"database/sql"
	"time"
)

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
