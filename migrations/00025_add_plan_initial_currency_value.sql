-- +goose Up
-- this is needed because some currency symbols are different accross exchanges
-- example: BCH is BCC on binance
ALTER TABLE plans ADD COLUMN user_currency_balance_at_init decimal; 
ALTER TABLE plans RENAME COLUMN base_currency_symbol TO user_currency_symbol;

-- +goose Down
ALTER TABLE plans DROP COLUMN user_currency_balance_at_init; 
ALTER TABLE plans RENAME COLUMN user_currency_symbol TO base_currency_symbol;
