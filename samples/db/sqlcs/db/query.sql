-- name: GetUser :one
SELECT id, first_name, last_name, display_name, bio, created_at FROM users
WHERE id = $1 LIMIT 1;

-- name: FindUserFirstName :one
SELECT id, first_name, last_name, display_name, bio, created_at FROM users
WHERE first_name like '%' || sqlc.arg(first_name) || '%' LIMIT 1;

-- name: FindUserGtCreatedAt :many
SELECT id, first_name, last_name, display_name, bio, created_at FROM users
WHERE created_at > $1 LIMIT $2;

-- name: ListUsers :many
SELECT id, first_name, last_name, display_name, bio, created_at FROM users
ORDER BY display_name;

-- name: CreateUser :execresult
INSERT INTO users (first_name, last_name, display_name, bio)
VALUES ($1, $2, $3, $4);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
