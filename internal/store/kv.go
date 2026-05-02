package store

import "database/sql"

func GetVersion(db *sql.DB, key string) (int, string, error) {
	var version int
	var value string

	err := db.QueryRow(`
		SELECT version, value
		FROM kv
		WHERE key = ?
	`, key).Scan(&version, &value)

	if err == sql.ErrNoRows {
		return 0, "", nil
	}

	return version, value, err
}
