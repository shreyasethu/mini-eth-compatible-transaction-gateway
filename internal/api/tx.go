package api

import "database/sql"

func GetTransaction(db *sql.DB, hash string) (map[string]interface{}, error) {
	row := db.QueryRow(`
		SELECT sender, nonce, status, replaced_by
		FROM txs
		WHERE hash = ?
	`, hash)

	var sender string
	var nonce int
	var status string
	var replaced string

	err := row.Scan(&sender, &nonce, &status, &replaced)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"hash":       hash,
		"from":       sender,
		"nonce":      nonce,
		"status":     status,
		"replacedBy": replaced,
	}, nil
}