-- +goose Up
-- SQL in this section is executed when the migration is applied.
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

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user_keys;

