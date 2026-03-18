-- name: ListUsersForAdmin :many
SELECT id, name, email, password, avatar, gender, bio, created_at, updated_at, deleted_type, deleted_at
FROM users
WHERE (@name::text IS NULL OR name ILIKE '%' || @name || '%')  -- name参数有值则匹配
  AND (@email::text IS NULL OR email ILIKE '%' || @email || '%')  -- email参数有值则匹配
  AND (@gender::text IS NULL OR TRIM(gender) = @gender)  -- gender参数有值则匹配
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUsers :many
SELECT id, name, email, password, avatar, gender, bio, created_at, updated_at
FROM users
WHERE (@name::text IS NULL OR name ILIKE '%' || @name || '%')  -- name参数有值则匹配
  AND (@email::text IS NULL OR email ILIKE '%' || @email || '%')  -- email参数有值则匹配
  AND (@id::text IS NULL OR id = @id)  -- id参数有值则匹配
  AND deleted_at IS NULL AND deleted_type = 0
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :exec
INSERT INTO
    users (
        id,
        name,
        email,
        password,
        avatar,
        gender,
        bio,
        deleted_type,
        deleted_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        0,
        NULL
    );

-- name: UpdateUser :exec
UPDATE users
SET
    name = coalesce(sqlc.narg('name'), name),
    email = coalesce(sqlc.narg('email'), email),
    password = coalesce(sqlc.narg('password'), password),
    avatar = coalesce(sqlc.narg('avatar'), avatar),
    gender = coalesce(sqlc.narg('gender'), gender),
    bio = COALESCE($1, bio),
    updated_at = now()
WHERE
    id = $2
    AND deleted_at IS NULL
    AND deleted_type = 0 RETURNING id,
    name,
    email,
    password,
    avatar,
    gender,
    bio,
    updated_at;

-- name: SoftDeleteUser :exec
UPDATE users
SET
    deleted_type = 1,
    deleted_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND deleted_at IS NULL
    AND deleted_type = 0;

-- name: RestoreUser :exec
UPDATE users
SET
    deleted_type = 0,
    deleted_at = NULL,
    updated_at = now()
WHERE
    id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;