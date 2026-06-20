-- name: GetExpense :one
SELECT
    e.id,
    e.user_id,
    e.category_id,
    e.amount,
    e.description,
    e.expense_date,
    e.payment_method,
    e.memo,
    e.created_at,
    e.updated_at,
    c.name AS category_name
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.id = $1 AND e.user_id = $2;

-- name: ListExpenses :many
SELECT
    e.id,
    e.user_id,
    e.category_id,
    e.amount,
    e.description,
    e.expense_date,
    e.payment_method,
    e.memo,
    e.created_at,
    e.updated_at,
    c.name AS category_name
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1 AND e.expense_date <= CURRENT_DATE
ORDER BY e.expense_date DESC;

-- name: ListExpensesByDateRange :many
SELECT
    e.id,
    e.user_id,
    e.category_id,
    e.amount,
    e.description,
    e.expense_date,
    e.payment_method,
    e.memo,
    e.created_at,
    e.updated_at,
    c.name AS category_name
FROM expenses e
         LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1
  AND e.expense_date BETWEEN $2 AND $3
ORDER BY e.expense_date DESC;

-- name: ListPlannedExpenses :many
SELECT
    e.id,
    e.user_id,
    e.category_id,
    e.amount,
    e.description,
    e.expense_date,
    e.payment_method,
    e.memo,
    e.created_at,
    e.updated_at,
    c.name AS category_name
FROM expenses e
         LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1
  AND e.expense_date > CURRENT_DATE
ORDER BY e.expense_date ASC;

-- name: SearchExpenses :many
SELECT
    e.id,
    e.user_id,
    e.category_id,
    e.amount,
    e.description,
    e.expense_date,
    e.payment_method,
    e.memo,
    e.created_at,
    e.updated_at,
    c.name AS category_name
FROM expenses e
         LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1
  AND (sqlc.narg('date_from')::date IS NULL OR e.expense_date >= sqlc.narg('date_from'))
  AND (sqlc.narg('date_to')::date IS NULL OR e.expense_date <= sqlc.narg('date_to'))
  AND (sqlc.narg('category_id')::uuid IS NULL OR e.category_id = sqlc.narg('category_id'))
  AND (sqlc.narg('keyword')::text IS NULL OR e.description ILIKE '%' || sqlc.narg('keyword') || '%')
ORDER BY e.expense_date DESC;

-- name: CreateExpense :one
INSERT INTO expenses (
    user_id,
    category_id,
    amount,
    description,
    expense_date,
    payment_method,
    memo
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         )
    RETURNING *;

-- name: UpdateExpense :one
UPDATE expenses SET
                    category_id    = $1,
                    amount         = $2,
                    description    = $3,
                    expense_date   = $4,
                    payment_method = $5,
                    memo           = $6,
                    updated_at     = CURRENT_TIMESTAMP
WHERE id = $7 AND user_id = $8
    RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses
WHERE id = $1 AND user_id = $2;
