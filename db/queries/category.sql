-- name: ListCategories :many
SELECT
    id,
    user_id,
    name,
    type,
    color,
    sort_order,
    is_default,
    created_at
FROM categories
WHERE user_id = $1
ORDER BY sort_order ASC;

-- name: GetCategory :one
SELECT
    id,
    user_id,
    name,
    type,
    color,
    sort_order,
    is_default,
    created_at
FROM categories
WHERE id = $1 AND user_id = $2;

-- name: CreateCategory :one
INSERT INTO categories (
    user_id,
    name,
    type,
    color,
    sort_order,
    is_default
) VALUES (
             $1, $2, $3, $4, $5, $6
         )
    RETURNING *;

-- name: UpdateCategory :one
UPDATE categories SET
                      name       = $1,
                      type       = $2,
                      color      = $3,
                      sort_order = $4
WHERE id = $5 AND user_id = $6
    RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1 AND user_id = $2;