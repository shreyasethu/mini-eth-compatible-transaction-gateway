# Mini Ethereum Compatible Transaction Gateway

## Lifecycle

RECEIVED -> PENDING -> INFLIGHT -> COMMITTED

or

PENDING -> FAILED

## Restart Recovery

On startup:

INFLIGHT -> PENDING

## Receipt Rules

Only COMMITTED tx gets receipt.

## Nonce Rule

Returns committed_nonce + 1