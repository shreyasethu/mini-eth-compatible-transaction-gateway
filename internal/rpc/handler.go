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
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Error:   "failed to read request body",
				ID:      0,
			})
			return
		}

		var req Request
		err = json.Unmarshal(body, &req)
		if err != nil {
			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Error:   "invalid json",
				ID:      0,
			})
			return
		}

		switch req.Method {

		case "eth_sendRawTransaction":
			if len(req.Params) == 0 {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Error:   "missing params",
					ID:      req.ID,
				})
				return
			}

			raw := string(req.Params[0])

			hashBytes := sha256.Sum256([]byte(raw))
			hash := "0x" + hex.EncodeToString(hashBytes[:])

			var payload types.RawTx
			err = json.Unmarshal(req.Params[0], &payload)
			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Error:   "invalid raw tx",
					ID:      req.ID,
				})
				return
			}

			_, err = db.Exec(`
				UPDATE txs
				SET status='REPLACED', replaced_by=?
				WHERE sender=?
				AND nonce=?
				AND status!='COMMITTED'
			`, hash, payload.From, payload.Nonce)

			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Error:   err.Error(),
					ID:      req.ID,
				})
				return
			}

			_, err = db.Exec(`
				INSERT INTO txs(hash,sender,nonce,raw,status)
				VALUES(?,?,?,?,?)
			`, hash, payload.From, payload.Nonce, raw, "PENDING")

			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Error:   err.Error(),
					ID:      req.ID,
				})
				return
			}

			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Result:  hash,
				ID:      req.ID,
			})

		case "eth_getTransactionByHash":
			if len(req.Params) == 0 {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Result:  nil,
					ID:      req.ID,
				})
				return
			}

			var hash string
			err = json.Unmarshal(req.Params[0], &hash)
			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Result:  nil,
					ID:      req.ID,
				})
				return
			}

			var tx types.Transaction

			row := db.QueryRow(`
				SELECT hash, sender, nonce, raw, status
				FROM txs
				WHERE hash=?
			`, hash)

			err = row.Scan(
				&tx.Hash,
				&tx.Sender,
				&tx.Nonce,
				&tx.Raw,
				&tx.Status,
			)

			if err != nil {
				json.NewEncoder(w).Encode(Response{
					Jsonrpc: "2.0",
					Result:  nil,
					ID:      req.ID,
				})
				return
			}

			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Result:  tx,
				ID:      req.ID,
			})

		default:
			json.NewEncoder(w).Encode(Response{
				Jsonrpc: "2.0",
				Error:   "method not found",
				ID:      req.ID,
			})
		}
	}
}
