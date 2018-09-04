-- +goose Up
ALTER TABLE orders DROP COLUMN user_key_id; 
ALTER TABLE orders ADD COLUMN account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE CASCADE;


-- +goose Down
ALTER TABLE orders ADD COLUMN user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE;
ALTER TABLE orders DROP COLUMN account_id; 