-- ユーザー
INSERT INTO users (id, email, password_hash, provider, name) VALUES
    ('11111111-1111-1111-1111-111111111111', 'test@example.com', '$2a$10$dummy_hash', 'LOCAL', 'テストユーザー');

-- カテゴリ（支出）
INSERT INTO categories (id, user_id, name, type, color, sort_order, is_default) VALUES
                                                                                    ('aaaaaaaa-0001-0001-0001-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', '食費',       'EXPENSE', '#FF6384', 1, TRUE),
                                                                                    ('aaaaaaaa-0002-0002-0002-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', '交通費',     'EXPENSE', '#36A2EB', 2, TRUE),
                                                                                    ('aaaaaaaa-0003-0003-0003-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', '日用品',     'EXPENSE', '#FFCE56', 3, TRUE),
                                                                                    ('aaaaaaaa-0004-0004-0004-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', '娯楽',       'EXPENSE', '#4BC0C0', 4, TRUE),
                                                                                    ('aaaaaaaa-0005-0005-0005-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', '光熱費',     'EXPENSE', '#9966FF', 5, TRUE);

-- カテゴリ（収入）
INSERT INTO categories (id, user_id, name, type, color, sort_order, is_default) VALUES
                                                                                    ('bbbbbbbb-0001-0001-0001-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', '給与',       'INCOME', '#4CAF50', 1, TRUE),
                                                                                    ('bbbbbbbb-0002-0002-0002-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', '副業',       'INCOME', '#8BC34A', 2, TRUE),
                                                                                    ('bbbbbbbb-0003-0003-0003-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'その他収入', 'INCOME', '#CDDC39', 3, TRUE);

-- 支出
INSERT INTO expenses (user_id, category_id, amount, description, expense_date, payment_method, memo) VALUES
                                                                                                         ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-0001-0001-0001-aaaaaaaaaaaa', 1500.00, 'ランチ代',         '2026-04-14', 'CASH',        '同僚と'),
                                                                                                         ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-0001-0001-0001-aaaaaaaaaaaa', 3200.00, '夕食（外食）',     '2026-04-15', 'CREDIT_CARD', NULL),
                                                                                                         ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-0002-0002-0002-aaaaaaaaaaaa',  230.00, '電車代',           '2026-04-15', 'CASH',        NULL),
                                                                                                         ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-0003-0003-0003-aaaaaaaaaaaa', 1980.00, 'シャンプー等',     '2026-04-16', 'CREDIT_CARD', 'ドラッグストア'),
                                                                                                         ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-0004-0004-0004-aaaaaaaaaaaa', 2000.00, '映画チケット',     '2026-04-17', 'CREDIT_CARD', NULL),
                                                                                                         ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-0005-0005-0005-aaaaaaaaaaaa', 8500.00, '電気代',           '2026-04-18', 'CREDIT_CARD', '4月分');

-- 収入
INSERT INTO incomes (user_id, category_id, amount, description, income_date, memo) VALUES
                                                                                       ('11111111-1111-1111-1111-111111111111', 'bbbbbbbb-0001-0001-0001-bbbbbbbbbbbb', 280000.00, '4月給与',   '2026-04-25', NULL),
                                                                                       ('11111111-1111-1111-1111-111111111111', 'bbbbbbbb-0002-0002-0002-bbbbbbbbbbbb',  50000.00, 'フリーランス案件', '2026-04-20', 'A社');