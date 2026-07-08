-- name: GetUser :one
SELECT id, name, bio, created_at FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, name, bio, created_at FROM users
ORDER BY name;

-- name: CreateUser :execresult
INSERT INTO users (name, bio)
VALUES (?, ?);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
