-- name: ListComments :many
SELECT c.id, c.user_id, c.content, c.post_id, c.created_at, c.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY c.created_at DESC;

-- name: GetCommentByID :one
SELECT c.id, c.user_id, c.content, c.post_id, c.created_at, c.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.id = ? AND c.deleted_at IS NULL AND u.deleted_at IS NULL
LIMIT 1;

-- name: ListCommentsByUserID :many
SELECT c.id, c.user_id, c.content, c.post_id, c.created_at, c.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.user_id = ? AND c.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY c.created_at DESC;

-- name: ListCommentsByPostID :many
SELECT c.id, c.user_id, c.content, c.post_id, c.created_at, c.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.post_id = ? AND c.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY c.created_at DESC;

-- name: ListCommentsByContent :many
SELECT c.id, c.user_id, c.content, c.post_id, c.created_at, c.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE INSTR(c.content, ?) > 0 AND c.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY c.created_at DESC;

-- name: CreateComment :exec
INSERT INTO comments (id, user_id, content, post_id)
VALUES (?, ?, ?, ?);

-- name: UpdateComment :exec
UPDATE comments
SET user_id = ?, content = ?, post_id = ?, updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteComment :exec
UPDATE comments
SET deleted_at = datetime('now', 'localtime'), updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: RestoreComment :exec
UPDATE comments
SET deleted_at = NULL, updated_at = datetime('now', 'localtime')
WHERE id = ?;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = ?;