-- name: UpsertUserByGoogleAuthSafe :one
INSERT INTO users (name, google_id, email, avatar_url, updated_at)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO UPDATE SET
    name = CASE
        WHEN users.name IS NULL OR users.name = '' THEN EXCLUDED.name
        ELSE COALESCE(EXCLUDED.name, users.name)
    END,
    google_id = COALESCE(users.google_id, EXCLUDED.google_id),
    avatar_url = COALESCE(EXCLUDED.avatar_url, users.avatar_url),
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, name, google_id, email, avatar_url, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByGoogleID :one
SELECT id, name, google_id, email, avatar_url, created_at, updated_at
FROM users
WHERE google_id = $1;


-- delete all except the one with email ending in @gmail.com
-- name: RemoveAllUsers :exec
DELETE FROM users
WHERE email NOT LIKE '%@gmail.com';