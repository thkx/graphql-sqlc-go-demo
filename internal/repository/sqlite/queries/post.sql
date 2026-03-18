-- name: ListPosts :many
SELECT p.id, p.title, p.content, p.pv, p.user_id, p.created_at, p.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY p.created_at DESC;

-- name: GetPostByID :one
SELECT p.id, p.title, p.content, p.pv, p.user_id, p.created_at, p.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.id = ? AND p.deleted_at IS NULL AND u.deleted_at IS NULL
LIMIT 1;

-- name: ListPostsByTitle :many
SELECT p.id, p.title, p.content, p.pv, p.user_id, p.created_at, p.updated_at,
       u.name AS user_name, u.email AS user_email, u.avatar AS user_avatar, 
       u.gender AS user_gender, u.bio AS user_bio
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE INSTR(p.title, ?) > 0 AND p.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY p.created_at DESC;

-- name: CreatePost :exec
INSERT INTO posts (id, title, content, pv, user_id)
VALUES (?, ?, ?, ?, ?);

-- name: UpdatePost :exec
UPDATE posts
SET title = ?, content = ?, pv = ?, user_id = ?, updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeletePost :exec
UPDATE posts
SET deleted_at = datetime('now', 'localtime'), updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: RestorePost :exec
UPDATE posts
SET deleted_at = NULL, updated_at = datetime('now', 'localtime')
WHERE id = ?;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = ?;