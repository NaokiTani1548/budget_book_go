ALTER TABLE expenses
DROP COLUMN IF EXISTS is_planned,
    DROP COLUMN IF EXISTS planned_date;

ALTER TABLE incomes
DROP COLUMN IF EXISTS is_planned,
    DROP COLUMN IF EXISTS planned_date;