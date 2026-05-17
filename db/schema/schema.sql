CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255),
                       provider VARCHAR(50) NOT NULL DEFAULT 'LOCAL',
                       provider_id VARCHAR(255),
                       name VARCHAR(100),
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                            name VARCHAR(100) NOT NULL,
                            type VARCHAR(20) NOT NULL,
                            color VARCHAR(7),
                            sort_order INTEGER NOT NULL DEFAULT 0,
                            is_default BOOLEAN NOT NULL DEFAULT FALSE,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            CONSTRAINT uk_categories_user_name_type UNIQUE (user_id, name, type),
                            CONSTRAINT chk_categories_type CHECK (type IN ('EXPENSE', 'INCOME'))
);

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

CREATE TABLE recurring_expenses (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
                                    amount DECIMAL(12, 2) NOT NULL,
                                    description VARCHAR(500),
                                    payment_method VARCHAR(50),
                                    memo TEXT,
                                    billing_day INTEGER NOT NULL,
                                    start_date DATE NOT NULL,
                                    end_date DATE,
                                    is_active BOOLEAN NOT NULL DEFAULT TRUE,
                                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    CONSTRAINT chk_recurring_billing_day CHECK (billing_day BETWEEN 1 AND 31),
                                    CONSTRAINT chk_recurring_amount CHECK (amount > 0),
                                    CONSTRAINT chk_recurring_payment_method CHECK (
                                        payment_method IS NULL OR payment_method IN ('CASH', 'CREDIT_CARD')
                                        )
);

CREATE TABLE recurring_expense_logs (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                        recurring_expense_id UUID NOT NULL REFERENCES recurring_expenses(id) ON DELETE CASCADE,
                                        expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
                                        billing_year INTEGER NOT NULL,
                                        billing_month INTEGER NOT NULL,
                                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT uk_recurring_expense_logs UNIQUE (recurring_expense_id, billing_year, billing_month)
);