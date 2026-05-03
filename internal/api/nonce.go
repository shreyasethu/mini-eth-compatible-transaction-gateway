package api

import "database/sql"

func GetTransactionCount(db *sql.DB, sender string) (int, error) {
	row := db.QueryRow(`
		SELECT committed_nonce
		FROM accounts
		WHERE sender = ?
	`, sender)

	var nonce int
	err := row.Scan(&nonce)

	if err == sql.ErrNoRows {
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	return nonce + 1, nil
}