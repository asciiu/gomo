-- +goose Up
CREATE TABLE trades (
  id UUID PRIMARY KEY NOT NULL,
  order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  initial_currency_symbol text NOT NULL,
  initial_currency_balance decimal DEFAULT 0,
  initial_currency_traded decimal DEFAULT 0,
  initial_currency_remainder decimal DEFAULT 0,
  initial_currency_price decimal DEFAULT 0,
  final_currency_symbol text NOT NULL,       
  final_currency_balance decimal DEFAULT 0,   
  fee_currency_symbol text,
  fee_currency_amount decimal DEFAULT 0,
  exchange_time TIMESTAMP, 
  side text NOT NULL,                 -- buy, sell
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);


-- +goose Down
DROP TABLE trades; 
