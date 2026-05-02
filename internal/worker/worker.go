package worker

import (
	"database/sql"
	"time"

	"mini-eth-compatible-transaction-gateway/internal/executor"
	"mini-eth-compatible-transaction-gateway/internal/store"
	"mini-eth-compatible-transaction-gateway/internal/util"
)

func Start(db *sql.DB) {
	go func() {
		for {
			RunOnce(db)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func RunOnce(db *sql.DB) {
	hash, sender, nonce, err := store.ClaimNext(db)
	if err != nil {
		return
	}

	// confirm tx is still valid after claim
	var status string
	err = db.QueryRow(`
		SELECT status
		FROM txs
		WHERE hash=?
	`, hash).Scan(&status)

	if err != nil {
		return
	}

	// replaced / changed tx should not continue
	if status != "INFLIGHT" {
		return
	}

	committed, err := store.GetCommittedNonce(db, sender)
	if err != nil {
		return
	}

	expected := committed + 1

	if nonce < expected {
		db.Exec(`
			UPDATE txs
			SET status='FAILED',
			    last_error='nonce too low'
			WHERE hash=?
			AND status='INFLIGHT'
		`, hash)
		return
	}

	if nonce > expected {
		db.Exec(`
			UPDATE txs
			SET status='PARKED'
			WHERE hash=?
			AND status='INFLIGHT'
		`, hash)
		return
	}

	rw, err := executor.Execute(db, "shared", hash)
	if err != nil {
		db.Exec(`
			UPDATE txs
			SET status='PENDING',
			    last_error='execute failed'
			WHERE hash=?
			AND status='INFLIGHT'
		`, hash)
		return
	}

	Commit(db, hash, sender, rw)
}

func Commit(db *sql.DB, hash string, sender string, rw executor.RWSet) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	// ensure still inflight
	var status string
	err = tx.QueryRow(`
		SELECT status
		FROM txs
		WHERE hash=?
	`, hash).Scan(&status)

	if err != nil {
		tx.Rollback()
		return
	}

	if status != "INFLIGHT" {
		tx.Rollback()
		return
	}

	var currentVersion int

	err = tx.QueryRow(`
		SELECT version
		FROM kv
		WHERE key=?
	`, rw.Key).Scan(&currentVersion)

	if err != nil {
		currentVersion = 0
	}

	// MVCC conflict
	if currentVersion != rw.ReadVersion {
		tx.Rollback()

		var retries int
		db.QueryRow(`
			SELECT retries
			FROM txs
			WHERE hash=?
		`, hash).Scan(&retries)

		next := time.Now().Add(util.RetryDelay(retries))

		db.Exec(`
			UPDATE txs
			SET status='PENDING',
			    retries=retries+1,
			    next_retry_at=?,
			    last_error='mvcc conflict'
			WHERE hash=?
			AND status='INFLIGHT'
		`, next, hash)

		return
	}

	_, err = tx.Exec(`
		INSERT OR REPLACE INTO kv(key,value,version)
		VALUES(?,?,?)
	`, rw.Key, rw.NewValue, currentVersion+1)

	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec(`
		INSERT OR REPLACE INTO receipts(hash,status,logs_json,created_at)
		VALUES(?, 'SUCCESS', '[]', CURRENT_TIMESTAMP)
	`, hash)

	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec(`
		UPDATE txs
		SET status='COMMITTED',
		    last_error=''
		WHERE hash=?
		AND status='INFLIGHT'
	`, hash)

	if err != nil {
		tx.Rollback()
		return
	}

	err = store.IncrementNonce(tx, sender)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	// unpark next nonce
	db.Exec(`
		UPDATE txs
		SET status='PENDING'
		WHERE sender=?
		AND nonce=(
			SELECT committed_nonce + 1
			FROM accounts
			WHERE sender=?
		)
		AND status='PARKED'
	`, sender, sender)
}
