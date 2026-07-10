-- name: GetUser :one
SELECT id, first_name, last_name, display_name, bio, created_at FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, first_name, last_name, display_name, bio, created_at FROM users
ORDER BY display_name;

-- name: CreateUser :execresult
INSERT INTO users (first_name, last_name, display_name, bio)
VALUES ($1, $2, $3, $4);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
