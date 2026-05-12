-- name: GetIncome :one
SELECT
    i.id,
    i.user_id,
    i.category_id,
    i.amount,
    i.description,
    i.income_date,
    i.memo,
    i.is_planned,
    i.planned_date,
    i.created_at,
    i.updated_at,
    c.name AS category_name
FROM incomes i
         LEFT JOIN categories c ON i.category_id = c.id
WHERE i.id = $1 AND i.user_id = $2;

-- name: ListIncomes :many
SELECT
    i.id,
    i.user_id,
    i.category_id,
    i.amount,
    i.description,
    i.income_date,
    i.memo,
    i.is_planned,
    i.planned_date,
    i.created_at,
    i.updated_at,
    c.name AS category_name
FROM incomes i
         LEFT JOIN categories c ON i.category_id = c.id
WHERE i.user_id = $1 AND i.is_planned = FALSE
ORDER BY i.income_date DESC;

-- name: ListPlannedIncomes :many
SELECT
    i.id,
    i.user_id,
    i.category_id,
    i.amount,
    i.description,
    i.income_date,
    i.memo,
    i.is_planned,
    i.planned_date,
    i.created_at,
    i.updated_at,
    c.name AS category_name
FROM incomes i
         LEFT JOIN categories c ON i.category_id = c.id
WHERE i.user_id = $1 AND i.is_planned = TRUE AND i.planned_date > CURRENT_DATE
ORDER BY i.planned_date ASC;

-- name: ListIncomesByDateRange :many
SELECT
    i.id,
    i.user_id,
    i.category_id,
    i.amount,
    i.description,
    i.income_date,
    i.memo,
    i.is_planned,
    i.planned_date,
    i.created_at,
    i.updated_at,
    c.name AS category_name
FROM incomes i
         LEFT JOIN categories c ON i.category_id = c.id
WHERE i.user_id = $1
  AND i.is_planned = FALSE
  AND i.income_date BETWEEN $2 AND $3
ORDER BY i.income_date DESC;

-- name: CreateIncome :one
INSERT INTO incomes (
    user_id,
    category_id,
    amount,
    description,
    income_date,
    memo,
    is_planned,
    planned_date
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8
         )
    RETURNING *;

-- name: UpdateIncome :one
UPDATE incomes SET
                   category_id  = $1,
                   amount       = $2,
                   description  = $3,
                   income_date  = $4,
                   memo         = $5,
                   is_planned   = $6,
                   planned_date = $7,
                   updated_at   = CURRENT_TIMESTAMP
WHERE id = $8 AND user_id = $9
    RETURNING *;

-- name: DeleteIncome :exec
DELETE FROM incomes
WHERE id = $1 AND user_id = $2;