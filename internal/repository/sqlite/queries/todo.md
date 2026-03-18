
```sql
-- name: ListTodos :many
SELECT t.id, t.text, t.done, t.user_id, t.created_at, t.updated_at,
       u.name AS user_name, u.email AS user_email
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.deleted_at IS NULL AND u.deleted_at IS NULL
ORDER BY t.created_at DESC;

-- name: GetTodoByID :one
SELECT t.id, t.text, t.done, t.user_id, t.created_at, t.updated_at,
       u.name AS user_name, u.email AS user_email
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.id = ? AND t.deleted_at IS NULL AND u.deleted_at IS NULL
LIMIT 1;

-- name: CreateTodo :exec
INSERT INTO todos (id, text, done, user_id)
VALUES (?, ?, ?, ?);

-- name: UpdateTodo :exec
UPDATE todos
SET text = ?, done = ?, updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteTodo :exec
UPDATE todos
SET deleted_at = datetime('now', 'localtime'), updated_at = datetime('now', 'localtime')
WHERE id = ? AND deleted_at IS NULL;

-- name: ListTodosByUserID :many
SELECT t.id, t.text, t.done, t.user_id, t.created_at, t.updated_at
FROM todos t
WHERE t.user_id = ? AND t.deleted_at IS NULL
ORDER BY t.created_at DESC;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;
```
