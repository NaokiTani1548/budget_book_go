CREATE TABLE incomes (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
                         amount DECIMAL(12, 2) NOT NULL,
                         description VARCHAR(500),
                         income_date DATE NOT NULL,
                         memo TEXT,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         CONSTRAINT chk_incomes_amount CHECK (amount > 0)
);