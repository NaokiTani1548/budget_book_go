-- name: SumActualIncomes :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM incomes
WHERE user_id = $1 AND is_planned = FALSE;

-- name: SumActualExpenses :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM expenses
WHERE user_id = $1 AND is_planned = FALSE;

-- name: SumPlannedIncomesByDate :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM incomes
WHERE user_id = $1
  AND is_planned = TRUE
  AND planned_date <= $2;

-- name: SumPlannedExpensesByDate :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM expenses
WHERE user_id = $1
  AND is_planned = TRUE
  AND planned_date <= $2;