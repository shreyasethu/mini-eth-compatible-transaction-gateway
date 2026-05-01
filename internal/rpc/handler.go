package rpc

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"mini-eth-compatible-transaction-gateway/internal/types"
)

type Request struct {
	Jsonrpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      int               `json:"id"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      int         `json:"id"`
}

func Handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		var req Request
		json.Unmarshal(body, &req)

		switch req.Method {

		case "eth_sendRawTransaction":
			raw := string(req.Params[0])

			hashBytes := sha256.Sum256([]byte(raw))
			hash := "0x" + hex.EncodeToString(hashBytes[:])

			_, err := db.Exec(
				"INSERT INTO txs(hash,sender,nonce,raw,status) VALUES(?,?,?,?,?)",
				hash, "unknown", 0, raw, "PENDING",
			)

			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Error: err.Error(),
					ID: req.ID,
				})
				return
			}

			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Result: hash,
				ID: req.ID,
			})

		case "eth_getTransactionByHash":
			var tx types.Transaction

			var hash string
			json.Unmarshal(req.Params[0], &hash)

			row := db.QueryRow(
				"SELECT hash,sender,nonce,raw,status FROM txs WHERE hash=?",
				hash,
			)

			err := row.Scan(
				&tx.Hash,
				&tx.Sender,
				&tx.Nonce,
				&tx.Raw,
				&tx.Status,
			)

			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Result: nil,
					ID: req.ID,
				})
				return
			}

			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Result: tx,
				ID: req.ID,
			})
		}
	}
}