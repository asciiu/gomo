-- +goose Up
-- this is needed because some currency symbols are different accross exchanges
-- example: BCH is BCC on binance
ALTER TABLE plans ADD COLUMN reference_price decimal; 

-- +goose Down
ALTER TABLE plans DROP COLUMN reference_price; 
