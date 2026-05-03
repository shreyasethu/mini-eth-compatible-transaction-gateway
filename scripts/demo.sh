#!/bin/bash

curl -X POST http://localhost:8080 \
-H "Content-Type: application/json" \
-d '{
  "jsonrpc":"2.0",
  "method":"eth_sendRawTransaction",
  "params":[
    {
      "from":"0xabc",
      "nonce":1
    }
  ],
  "id":1
}'

echo

curl -X POST http://localhost:8080 \
-H "Content-Type: application/json" \
-d '{
  "jsonrpc":"2.0",
  "method":"eth_getTransactionCount",
  "params":["0xabc","latest"],
  "id":2
}'