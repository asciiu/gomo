-- +goose Up
ALTER TABLE orders DROP COLUMN user_key_id; 
ALTER TABLE orders ADD COLUMN account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE CASCADE;
DROP TABLE user_balances;
DROP TABLE user_keys;


-- +goose Down
ALTER TABLE orders ADD COLUMN user_key_id UUID NOT NULL REFERENCES user_keys (id) ON DELETE CASCADE;
ALTER TABLE orders DROP COLUMN account_id; 

CREATE TABLE user_keys(
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,
  exchange_name VARCHAR NOT NULL,
  api_key VARCHAR NOT NULL,
  secret  VARCHAR NOT NULL,
  description VARCHAR,
  status VARCHAR NOT NULL,
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP default current_timestamp,
  UNIQUE (api_key, secret)
);

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