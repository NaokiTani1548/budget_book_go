-- name: SumActualIncomes :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM incomes
WHERE user_id = $1
  AND income_date <= CURRENT_DATE;

-- name: SumActualExpenses :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM expenses
WHERE user_id = $1
  AND expense_date <= CURRENT_DATE;

-- name: SumPlannedIncomesByDate :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM incomes
WHERE user_id = $1
  AND income_date > CURRENT_DATE
  AND income_date <= $2;

-- name: SumPlannedExpensesByDate :one
SELECT COALESCE(SUM(amount), 0)::NUMERIC AS total
FROM expenses
WHERE user_id = $1
  AND expense_date > CURRENT_DATE
  AND expense_date <= $2;