CREATE TABLE expenses (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
                          amount DECIMAL(12, 2) NOT NULL,
                          description VARCHAR(500),
                          expense_date DATE NOT NULL,
                          payment_method VARCHAR(50),
                          memo TEXT,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT chk_expenses_amount CHECK (amount > 0),
                          CONSTRAINT chk_expenses_payment_method CHECK (
                              payment_method IS NULL OR payment_method IN ('CASH', 'CREDIT_CARD')
                              )
);