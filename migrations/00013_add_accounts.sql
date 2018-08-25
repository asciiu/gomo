-- +goose Up
CREATE TABLE accounts (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  exchange_name text,
  key_public text,
  key_secret text, 
  description text,
  status text,
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE balances (
 id UUID PRIMARY KEY NOT NULL,
 user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
 account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
 currency_symbol text NOT NULL,
 blockchain_address text,
 available decimal NOT NULL default 0.0,           -- what is free to use in our system
 locked decimal NOT NULL default 0.0,              -- amount in use in our system
 exchange_total decimal NOT NULL default 0.0,      -- total as seen on exchange
 exchange_available decimal NOT NULL default 0.0,  -- free to use on exchange 
 exchange_locked decimal NOT NULL default 0.0,     -- currently in use on exchange
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP default current_timestamp,
 UNIQUE (account_id, currency_symbol)              -- no dupe currency holdings per account
);


-- +goose Down
DROP TABLE balances;
DROP TABLE accounts;

