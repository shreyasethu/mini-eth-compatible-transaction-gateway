package store

import (
	"database/sql"
	"time"
)

func ClaimNext(db *sql.DB) (string, string, uint64, error) {
	row := db.QueryRow(`
		SELECT hash, sender, nonce
		FROM txs
		WHERE status = 'PENDING'
		AND next_retry_at <= ?
		ORDER BY created_at
		LIMIT 1
	`, time.Now())

	var hash string
	var sender string
	var nonce uint64

	err := row.Scan(&hash, &sender, &nonce)
	if err != nil {
		return "", "", 0, err
	}

	_, err = db.Exec(`
		UPDATE txs
		SET status='INFLIGHT'
		WHERE hash=?
	`, hash)

	return hash, sender, nonce, err
}
