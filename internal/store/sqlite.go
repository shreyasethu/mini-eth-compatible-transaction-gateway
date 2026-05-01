package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data/gateway.db")
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS txs (
		hash TEXT PRIMARY KEY,
		sender TEXT,
		nonce INTEGER,
		raw TEXT,
		status TEXT
	);
	`

	_, err = db.Exec(query)
	return db, err
}