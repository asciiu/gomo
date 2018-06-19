-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE strategies (
 id UUID PRIMARY KEY NOT NULL,
 user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
 user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE,
 exchange_name text NOT NULL,
 market_name text NOT NULL,
 order_ids uuid[] NOT NULL,
 base_balance decimal NOT NULL,
 currency_balance decimal NOT NULL,
 status text NOT NULL,               -- strategy status is active, inactive, or failed
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE strategies;
