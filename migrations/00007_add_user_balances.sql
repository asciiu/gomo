-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE user_balances (
 id UUID PRIMARY KEY NOT NULL,
 user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
 user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE,
 exchange_name text NOT NULL,
 currency_name text NOT NULL,
 blockchain_address text,
 available decimal NOT NULL default 0.0,
 locked decimal NOT NULL default 0.0,
 exchange_total decimal NOT NULL default 0.0,
 exchange_available decimal NOT NULL default 0.0,
 exchange_locked decimal NOT NULL default 0.0,
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP default current_timestamp,
 UNIQUE (user_key_id, currency_name)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user_balances;