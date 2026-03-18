-- name: ListPostsForAdmin :many
SELECT
	p.id,
  p.title,
  p.content,
  p.pv,
  p.created_at,
  p.updated_at,
  p.deleted_type,
  p.deleted_at,
  u.id AS user_id,
  u.name AS user_name,
  u.email AS user_email,
  u.avatar AS user_avatar,
  u.gender AS user_gender,
  u.bio AS user_bio,
  u.deleted_type AS user_deleted_type,
  u.deleted_at AS user_deleted_at
FROM posts p
JOIN users u ON p.author = u.id
WHERE (@title::text IS NULL OR p.title ILIKE '%' || @title || '%')  -- title有值则按title筛选
  AND (@content::text IS NULL OR p.content ILIKE '%' || @content || '%')  -- content有值则按content筛选
  AND (@pv::integer IS NULL OR p.pv = @pv)  -- pv有值则按pv筛选
  AND (@author::text IS NULL OR p.author = @author)  -- author有值则匹配用户
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPosts :many
SELECT 
  p.id,
  p.title,
  p.content,
  p.pv,
  p.created_at,
  p.updated_at,
  u.id AS user_id,
  u.name AS user_name,
  u.email AS user_email,
  u.avatar AS user_avatar,
  u.gender AS user_gender,
  u.bio AS user_bio
FROM posts p
JOIN users u ON p.author = u.id
WHERE 
  (@title::text IS NULL OR p.title ILIKE '%' || @title || '%')  -- title有值则匹配标题
  AND (@content::text IS NULL OR p.content ILIKE '%' || @content || '%')  -- content有值则匹配内容
  AND (@pv::integer IS NULL OR p.pv = @pv)  -- pv有值则匹配浏览量
  AND (@author::text IS NULL OR p.author = @author)  -- author有值则匹配用户
  AND (@id::text IS NULL OR p.id = @id)  -- id参数有值则匹配
  AND p.deleted_at IS NULL
  AND p.deleted_type = 0
  AND u.deleted_at IS NULL
  AND u.deleted_type = 0
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreatePost :exec
INSERT INTO
    posts (
        id,
        title,
        content,
        pv,
        author,
        deleted_type,
        deleted_at
    )
VALUES ($1, $2, $3, $4, $5, 0, NULL);

-- name: UpdatePost :exec
UPDATE posts
SET
    title = coalesce(sqlc.narg('title'), title),
    content = coalesce(sqlc.narg('content'), content),
    author = coalesce(sqlc.narg('author'), author),
    updated_at = now()
WHERE
    id = $1
    AND deleted_at IS NULL
    AND deleted_type = 0 RETURNING id,
    title,
    content,
    author,
    updated_at;

-- name: SoftDeletePost :exec
UPDATE posts
SET
    deleted_type = 1,
    deleted_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND deleted_at IS NULL
    AND deleted_type = 0;

-- name: RestorePost :exec
UPDATE posts
SET
    deleted_type = 0,
    deleted_at = NULL,
    updated_at = now()
WHERE
    id = $1;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;