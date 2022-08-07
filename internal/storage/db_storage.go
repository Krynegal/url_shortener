package storage

import "database/sql"

type DB struct {
	db *sql.DB
}

func NewDB(database *sql.DB) *DB {
	db := &DB{
		db: database,
	}
	return db
}

func (db *DB) Shorten(uid string, url string) (int, error) {
	_, err := db.db.Exec("INSERT INTO URLS (original, created_by) VALUES ($1, $2);", url, uid)
	if err != nil {
		return -1, err
	}
	row := db.db.QueryRow("SELECT url_id FROM URLS WHERE original = ($1)", url)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) Unshorten(id string) (string, error) {
	row := db.db.QueryRow("SELECT original FROM URLS WHERE url_id = ($1)", id)
	var origURL string
	if err := row.Scan(&origURL); err != nil {
		return "", err
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
	}()

	for rows.Next() {
		var url, orig string
		if err := rows.Scan(&url, &orig); err != nil {
			return nil
		}
		allURLs[url] = orig
	}

	return allURLs
}
