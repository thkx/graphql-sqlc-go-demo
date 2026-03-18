-- name: ListIndicators :many
SELECT id, indicator, indicator_type, meta_source, created_at, updated_at
FROM indicators
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetIndicatorByID :one
SELECT id, indicator, indicator_type, meta_source, created_at, updated_at
FROM indicators
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: CreateIndicator :exec
INSERT INTO indicators (id, indicator, indicator_type, meta_source)
VALUES (?, ?, ?, ?);

-- name: UpdateIndicator :exec
UPDATE indicators
SET indicator = ?, indicator_type = ?, meta_source = ?, updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteIndicator :exec
UPDATE indicators
SET deleted_at = datetime('now', 'localtime'), updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: SearchIndicatorsByType :many
SELECT id, indicator, indicator_type, meta_source, created_at, updated_at
FROM indicators
WHERE indicator_type = ? AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteIndicator :exec
DELETE FROM indicators WHERE id = ?;
