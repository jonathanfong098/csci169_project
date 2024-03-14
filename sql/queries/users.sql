-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, api_key, email, subscribed)
VALUES (
    $1,
    $2,
    $3,
    $4,
    encode(sha256(random()::text::bytea), 'hex'),
    $5,
    FALSE
)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetSubscribedUsers :many
SELECT * FROM users WHERE subscribed = TRUE;

-- name: SubscribeUser :one
UPDATE users
SET subscribed = TRUE
WHERE id = $1
RETURNING *;

-- name: UnsubscribeUser :one
UPDATE users
SET subscribed = FALSE
WHERE id = $1
RETURNING *;
