-- name: ListComments :many
SELECT
  c.id,
  c.content,
  c.post_id,
  c.created_at,
  c.updated_at,
  u.id AS user_id,
  u.name AS user_name,
  u.email AS user_email,
  u.avatar AS user_avatar,
  u.gender AS user_gender,
  u.bio AS user_bio
FROM comments c
JOIN users u ON c.author = u.id
WHERE (@content::text IS NULL OR c.content ILIKE '%' || @content || '%')  -- content参数有值则匹配
  AND (@author::text IS NULL OR c.author = @author)  -- author参数有值则匹配
	AND (@post_id::text IS NULL OR c.post_id = @post_id)  -- post_id参数有值则匹配
  AND (@id::text IS NULL OR c.id = @id)  -- id参数有值则匹配
  AND c.deleted_at IS NULL
  AND u.deleted_at IS NULL
ORDER BY c.created_at DESC;


-- name: CreateComment :exec
INSERT INTO comments (id, content, author, post_id, deleted_type, deleted_at)
VALUES ($1, $2, $3, $4, 0, NULL);


-- name: UpdateComment :exec
UPDATE comments
SET
  content = coalesce(sqlc.narg('content'), content),
  updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, content, updated_at;


-- name: SoftDeleteComment :exec
UPDATE comments
SET
  deleted_at = now(),
  updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL;


-- name: RestoreComment :exec
UPDATE comments
SET
  deleted_type = 0,
  deleted_at = NULL,
  updated_at = now()
WHERE id = $1;


-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;
