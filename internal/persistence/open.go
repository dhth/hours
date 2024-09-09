package persistence

import (
	"database/sql"
)

func GetDB(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbpath)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, err
}
