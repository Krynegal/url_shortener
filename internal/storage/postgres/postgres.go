package postgres

import (
	"database/sql"
)

const (
	URLTable = `
	CREATE TABLE IF NOT EXISTS URLS
	(
		url_id  serial PRIMARY KEY,
		original text NOT NULL UNIQUE, 
		created_by text
	);
	`
)

func NewDatabaseStorage(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if _, err = db.Exec(URLTable); err != nil {
		return nil, err
	}
	return db, nil
}
