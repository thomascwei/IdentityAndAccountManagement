-- name: GetAccount :one
SELECT id, username, email, auth
FROM accounts
WHERE id = ?
LIMIT 1;

-- name: ListAccounts :many
SELECT id, username, email, auth
FROM accounts
ORDER BY id;

-- name: CreateAccount :execresult
INSERT INTO accounts (username, password, email, auth)
VALUES (?, ?, ?, ?);

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = ?;

-- name: UpdateAccount :exec
UPDATE accounts
SET username = ?,
    email    = ?,
    auth     = ?
WHERE id = ?;

-- name: UpdatePassword :exec
UPDATE accounts
SET password = ?
WHERE id = ?;
