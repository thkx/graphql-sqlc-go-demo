-- name: ListUsers :many
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetUserByID :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE email = ? AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, name, email, password)
VALUES (?, ?, ?, ?);

-- name: UpdateUser :exec
UPDATE users
SET name = ?, email = ?, password = ?, updated_at = now()
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = now(), updated_at = now()
WHERE id = ? AND deleted_at IS NULL;

-- name: RestoreUser :exec
UPDATE users
SET deleted_at = NULL, updated_at = now()
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;