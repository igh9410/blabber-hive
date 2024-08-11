-- +goose Up
-- +goose StatementBegin

SELECT 'ALTER TABLE ' || tc.table_name || ' DROP CONSTRAINT ' || tc.constraint_name || ';'
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
ON ccu.constraint_name = tc.constraint_name
WHERE constraint_type = 'FOREIGN KEY' AND tc.table_schema='public';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'ALTER TABLE ' || tc.table_name || ' ADD CONSTRAINT ' || tc.constraint_name || ' FOREIGN KEY (' || kcu.column_name || ') REFERENCES ' || ccu.table_name || '(' || ccu.column_name || ');'
-- +goose StatementEnd
