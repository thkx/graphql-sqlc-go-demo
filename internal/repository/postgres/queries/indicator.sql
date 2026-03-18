-- name: ListIndicators :many
SELECT id, indicator, indicator_type, meta_source, created_at, updated_at
FROM indicators
WHERE (@indicator::text IS NULL OR indicator ILIKE '%' || @indicator || '%')  -- indicator参数有值则匹配
  AND (@indicator_type::text IS NULL OR indicator_type ILIKE '%' || @indicator_type || '%')  -- indicator_type参数有值则匹配
  AND (@meta_source::text IS NULL OR meta_source ILIKE '%' || @meta_source || '%')  -- meta_source参数有值则匹配
  AND (@id::text IS NULL OR id = @id)  -- id参数有值则匹配
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: CreateIndicator :exec
INSERT INTO indicators (id, indicator, indicator_type, meta_source)
VALUES ($1, $2, $3, $4);


-- name: UpdateIndicator :exec
UPDATE indicators
SET
    indicator = coalesce(sqlc.narg('indicator'), indicator),
    indicator_type = coalesce(sqlc.narg('indicator_type'), indicator_type),
    meta_source = coalesce(sqlc.narg('meta_source'), meta_source),
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, indicator, indicator_type, meta_source, updated_at;


-- name: SoftDeleteIndicator :exec
UPDATE indicators
SET
    deleted_at = now(),
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL;


-- name: DeleteIndicator :exec
DELETE FROM indicators
WHERE id = $1;
