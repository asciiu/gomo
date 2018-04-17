-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE orders (
 id UUID PRIMARY KEY NOT NULL,
 user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
 user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE,
 exchange_name text NOT NULL,
 exchange_order_id text,
 exchange_market_name text,
 market_name text NOT NULL,
 side text NOT NULL,
 "type" text NOT NULL,
 price decimal,
 quantity decimal NOT NULL,
 quantity_remaining decimal NOT NULL,
 status text NOT NULL,
 conditions jsonb NOT NULL,
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE orders;
