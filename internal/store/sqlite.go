package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "gateway.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS txs (
		hash TEXT PRIMARY KEY,
		sender TEXT,
		nonce INTEGER,
		raw TEXT,
		status TEXT,
		replaced_by TEXT DEFAULT '',
		retries INTEGER DEFAULT 0,
		last_error TEXT DEFAULT '',
		next_retry_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS receipts (
		hash TEXT PRIMARY KEY,
		status TEXT,
		logs_json TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS accounts (
		sender TEXT PRIMARY KEY,
		committed_nonce INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS kv (
		key TEXT PRIMARY KEY,
		value TEXT,
		version INTEGER DEFAULT 0
	);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
