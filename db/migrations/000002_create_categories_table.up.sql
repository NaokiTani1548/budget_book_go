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