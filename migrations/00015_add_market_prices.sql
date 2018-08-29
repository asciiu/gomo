-- +goose Up
-- this shall be the currency used to measure initial and final currency valuations 
CREATE TABLE exchange_rates (
  exchange_name text NOT NULL,
  market_name text NOT NULL,
  closed_at_price decimal NOT NULL,
  closed_at_time TIMESTAMP
);

-- +goose Down
DROP TABLE exchange_rates;