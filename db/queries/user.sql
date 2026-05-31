-- name: GetUserByID :one
SELECT id, email, password_hash, provider, provider_id, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, provider, provider_id, name, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByProviderID :one
SELECT id, email, password_hash, provider, provider_id, name, created_at, updated_at
FROM users
WHERE provider = $1 AND provider_id = $2;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, provider, provider_id, name)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
                 name = $1,
                 updated_at = CURRENT_TIMESTAMP
WHERE id = $2
    RETURNING *;