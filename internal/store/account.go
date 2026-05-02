package store

import "database/sql"

func GetCommittedNonce(db *sql.DB, sender string) (uint64, error) {
	var nonce uint64

	err := db.QueryRow(`
		SELECT committed_nonce
		FROM accounts
		WHERE sender = ?
	`, sender).Scan(&nonce)

	if err == sql.ErrNoRows {
		_, err = db.Exec(`
			INSERT INTO accounts(sender, committed_nonce)
			VALUES(?, 0)
		`, sender)
		return 0, err
	}

	return nonce, err
}

func IncrementNonce(tx *sql.Tx, sender string) error {
	_, err := tx.Exec(`
		UPDATE accounts
		SET committed_nonce = committed_nonce + 1
		WHERE sender = ?
	`, sender)

	return err
}
