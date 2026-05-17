-- name: GetRecurringExpense :one
SELECT
    r.id,
    r.user_id,
    r.category_id,
    r.amount,
    r.description,
    r.payment_method,
    r.memo,
    r.billing_day,
    r.start_date,
    r.end_date,
    r.is_active,
    r.created_at,
    r.updated_at,
    c.name AS category_name
FROM recurring_expenses r
         LEFT JOIN categories c ON r.category_id = c.id
WHERE r.id = $1 AND r.user_id = $2;

-- name: ListRecurringExpenses :many
SELECT
    r.id,
    r.user_id,
    r.category_id,
    r.amount,
    r.description,
    r.payment_method,
    r.memo,
    r.billing_day,
    r.start_date,
    r.end_date,
    r.is_active,
    r.created_at,
    r.updated_at,
    c.name AS category_name
FROM recurring_expenses r
         LEFT JOIN categories c ON r.category_id = c.id
WHERE r.user_id = $1
ORDER BY r.billing_day ASC;

-- name: ListActiveRecurringExpenses :many
SELECT
    r.id,
    r.user_id,
    r.category_id,
    r.amount,
    r.description,
    r.payment_method,
    r.memo,
    r.billing_day,
    r.start_date,
    r.end_date,
    r.is_active,
    r.created_at,
    r.updated_at,
    c.name AS category_name
FROM recurring_expenses r
         LEFT JOIN categories c ON r.category_id = c.id
WHERE r.user_id = $1
  AND r.is_active = TRUE
  AND r.start_date <= CURRENT_DATE
  AND (r.end_date IS NULL OR r.end_date >= CURRENT_DATE)
ORDER BY r.billing_day ASC;

-- name: CreateRecurringExpense :one
INSERT INTO recurring_expenses (
    user_id,
    category_id,
    amount,
    description,
    payment_method,
    memo,
    billing_day,
    start_date,
    end_date
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         )
    RETURNING *;

-- name: UpdateRecurringExpense :one
UPDATE recurring_expenses SET
                              category_id    = $1,
                              amount         = $2,
                              description    = $3,
                              payment_method = $4,
                              memo           = $5,
                              billing_day    = $6,
                              start_date     = $7,
                              end_date       = $8,
                              is_active      = $9,
                              updated_at     = CURRENT_TIMESTAMP
WHERE id = $10 AND user_id = $11
    RETURNING *;

-- name: DeleteRecurringExpense :exec
DELETE FROM recurring_expenses
WHERE id = $1 AND user_id = $2;

-- name: GetRecurringExpenseLog :one
SELECT id FROM recurring_expense_logs
WHERE recurring_expense_id = $1
  AND billing_year = $2
  AND billing_month = $3;

-- name: CreateRecurringExpenseLog :one
INSERT INTO recurring_expense_logs (
    recurring_expense_id,
    expense_id,
    billing_year,
    billing_month
) VALUES (
             $1, $2, $3, $4
         )
    RETURNING *;