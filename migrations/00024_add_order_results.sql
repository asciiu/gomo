-- +goose Up
-- this is needed because some currency symbols are different accross exchanges
-- example: BCH is BCC on binance
ALTER TABLE orders ADD COLUMN fee_currency_symbol text;
ALTER TABLE orders ADD COLUMN fee_currency_amount decimal DEFAULT 0;
ALTER TABLE orders ADD COLUMN exchange_price decimal DEFAULT 0;
ALTER TABLE orders ADD COLUMN exchange_time TIMESTAMP;
ALTER TABLE orders ADD COLUMN errors text;

-- +goose Down
ALTER TABLE orders DROP COLUMN fee_currency_symbol; 
ALTER TABLE orders DROP COLUMN fee_currency_amount; 
ALTER TABLE orders DROP COLUMN exchange_price; 
ALTER TABLE orders DROP COLUMN exchange_time; 
ALTER TABLE orders DROP COLUMN errors; 
