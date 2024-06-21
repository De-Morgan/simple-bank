

-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListOutGoingTransfers :many
SELECT * FROM transfers
WHERE from_account_id = $1
LIMIT $2
OFFSET $3;

-- name: ListIncomingTransfers :many
SELECT * FROM transfers
WHERE to_account_id = $1
LIMIT $2
OFFSET $3;


-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;
