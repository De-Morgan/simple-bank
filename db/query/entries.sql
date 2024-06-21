
-- name: CreatEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListAccountEntries :many
SELECT * FROM entries
WHERE account_id = $1
LIMIT $2
OFFSET $3;


-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
