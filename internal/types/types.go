package types

import "time"

const (
	Pending   = "PENDING"
	Inflight  = "INFLIGHT"
	Committed = "COMMITTED"
	Failed    = "FAILED"
	Parked    = "PARKED"
	Replaced  = "REPLACED"
)

type RawTx struct {
	From   string            `json:"from"`
	Nonce  int               `json:"nonce"`
	To     string            `json:"to"`
	Gas    int               `json:"gas"`
	Fee    int               `json:"fee"`
	Data   string            `json:"data"`
	Writes map[string]string `json:"writes"`
	Reads  []string          `json:"reads"`
}

type Transaction struct {
	Hash   string `json:"hash"`
	Sender string `json:"sender"`
	Nonce  int    `json:"nonce"`
	Raw    string `json:"raw"`
	Status string `json:"status"`
}

type Tx struct {
	Hash        string
	Sender      string
	Nonce       uint64
	Raw         string
	Status      string
	ReplacedBy  string
	Retries     int
	NextRetryAt time.Time
	LastError   string
}

type Receipt struct {
	Hash   string
	Status string
	Logs   string
}
