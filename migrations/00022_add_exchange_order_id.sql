-- +goose Up
ALTER TABLE orders ADD COLUMN exchange_order_id text;

-- +goose Down
ALTER TABLE orders DROP COLUMN exchange_order_id; 
