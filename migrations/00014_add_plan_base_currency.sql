-- +goose Up
-- this shall be the currency used to measure initial and final currency valuations 
ALTER TABLE plans ADD COLUMN base_currency_symbol text; 

-- +goose Down
ALTER TABLE plans DROP COLUMN base_currency_symbol;