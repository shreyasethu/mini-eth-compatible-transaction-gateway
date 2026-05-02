package executor

import (
	"database/sql"
	"mini-eth-compatible-transaction-gateway/internal/store"
)

type RWSet struct {
	Key         string
	ReadVersion int
	NewValue    string
}

func Execute(db *sql.DB, key string, value string) (RWSet, error) {
	version, _, err := store.GetVersion(db, key)
	if err != nil {
		return RWSet{}, err
	}

	return RWSet{
		Key:         key,
		ReadVersion: version,
		NewValue:    value,
	}, nil
}
