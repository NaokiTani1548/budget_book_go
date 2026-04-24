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