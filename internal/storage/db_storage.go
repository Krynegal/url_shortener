package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
)

const (
	URLTable = `
	CREATE TABLE IF NOT EXISTS URLS
	(
		url_id  serial PRIMARY KEY,
		original text NOT NULL UNIQUE, 
		created_by text,
		deleted boolean NOT NULL DEFAULT false
	);
	`
)

type DBStorager interface {
	Storager
	Ping(ctx context.Context) error
	Delete(string, []int) error
}

type DB struct {
	db *sql.DB
}

func NewDatabaseStorage(dataSourceName string) (DBStorager, error) {
	database, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}
	if _, err = database.Exec(URLTable); err != nil {
		return nil, err
	}
	db := &DB{
		db: database,
	}
	return db, nil
}

func (db *DB) Ping(ctx context.Context) error {
	if err := db.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (db *DB) Delete(uid string, urlIDs []int) error {
	fmt.Println(urlIDs)
	var a pgtype.Int8Array
	a.Set(urlIDs)
	rows, _ := db.db.Query("UPDATE URLS SET deleted=true WHERE created_by = $1 AND url_id = ANY($2);", uid, a)
	err := rows.Close()
	if err != nil {
		return err
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Shorten(uid string, url string) (int, error) {
	var id int
	if url == "" {
		return -1, errors.New("url is empty")
	}
	_, err := db.db.Exec("INSERT INTO URLS (original, created_by) VALUES ($1, $2);", url, uid)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pgerrcode.UniqueViolation {
			row := db.db.QueryRow("SELECT url_id FROM URLS WHERE original = ($1)", url)
			if err := row.Scan(&id); err != nil {
				return -1, err
			}
			return id, ErrKeyExists
		}
		return -1, err
	}
	row := db.db.QueryRow("SELECT url_id FROM URLS WHERE original = ($1)", url)
	if err = row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) Unshorten(id string) (string, error) {
	row := db.db.QueryRow("SELECT original, deleted FROM URLS WHERE url_id = ($1)", id)
	var origURL string
	var deleted bool
	if err := row.Scan(&origURL, &deleted); err != nil {
		return "", err
	}
	if deleted {
		return "", ErrKeyDeleted
	}
	return origURL, nil
}

func (db *DB) GetAllURLs(uid string) map[string]string {
	allURLs := map[string]string{}
	rows, err := db.db.Query("SELECT url_id, original FROM URLS WHERE created_by = ($1)", uid)
	if err != nil {
		return nil
	}

	defer func() {
		cerr := rows.Close()
		if cerr != nil {
			err = cerr
		}
		_ = rows.Err()
	}()

	for rows.Next() {
		var url, orig string
		if err = rows.Scan(&url, &orig); err != nil {
			return nil
		}
		allURLs[url] = orig
	}
	return allURLs
}
