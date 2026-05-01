package main

import (
	"log"
	"net/http"

	"mini-eth-compatible-transaction-gateway/internal/rpc"
	"mini-eth-compatible-transaction-gateway/internal/store"
)

func main() {
	db, err := store.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rpc.Handler(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}