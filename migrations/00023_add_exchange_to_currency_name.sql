-- +goose Up
-- this is needed because some currency symbols are different accross exchanges
-- example: BCH is BCC on binance
ALTER TABLE currency_names ADD COLUMN exchange_name text DEFAULT '*';

-- +goose Down
ALTER TABLE currency_names DROP COLUMN exchange_name; 
