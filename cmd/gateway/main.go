package main

import (
	"log"
	"net/http"

	"mini-eth-compatible-transaction-gateway/internal/recovery"
	"mini-eth-compatible-transaction-gateway/internal/rpc"
	"mini-eth-compatible-transaction-gateway/internal/store"
	"mini-eth-compatible-transaction-gateway/internal/worker"
)

func main() {
	db, err := store.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = recovery.ResetInflight(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Recovered INFLIGHT transactions to PENDING")

	worker.Start(db)

	http.HandleFunc("/", rpc.Handler(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}