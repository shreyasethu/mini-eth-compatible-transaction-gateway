package recovery

import (
	"database/sql"
)

func ResetInflight(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE txs
		SET status='PENDING'
		WHERE status='INFLIGHT'
	`)
	return err
}