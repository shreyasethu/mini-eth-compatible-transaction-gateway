package store

import "database/sql"

func CommitTx(db *sql.DB, hash string, sender string, nonce int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE txs
		SET status='COMMITTED'
		WHERE hash=?
	`, hash)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT OR IGNORE INTO receipts(hash,status,block_number,logs_json)
		VALUES(?,1,1,'[]')
	`, hash)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO accounts(sender, committed_nonce)
		VALUES(?,?)
		ON CONFLICT(sender)
		DO UPDATE SET committed_nonce=excluded.committed_nonce
	`, sender, nonce)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}