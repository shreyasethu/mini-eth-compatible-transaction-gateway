package api

import "database/sql"

func GetReceipt(db *sql.DB, hash string) (map[string]interface{}, error) {
	row := db.QueryRow(`
		SELECT status, block_number, logs_json
		FROM receipts
		WHERE hash = ?
	`, hash)

	var status int
	var block int
	var logs string

	err := row.Scan(&status, &block, &logs)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transactionHash": hash,
		"status":          status,
		"blockNumber":     block,
		"logs":            logs,
	}, nil
}