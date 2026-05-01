package types

type RawTx struct {
	From  string            `json:"from"`
	Nonce int               `json:"nonce"`
	To    string            `json:"to"`
	Gas   int               `json:"gas"`
	Fee   int               `json:"fee"`
	Data  string            `json:"data"`
	Writes map[string]string `json:"writes"`
	Reads []string          `json:"reads"`
}

type Transaction struct {
	Hash   string `json:"hash"`
	Sender string `json:"sender"`
	Nonce  int    `json:"nonce"`
	Raw    string `json:"raw"`
	Status string `json:"status"`
}