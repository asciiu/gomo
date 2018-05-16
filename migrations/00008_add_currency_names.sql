-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE currency_names (
  currency_name text NOT NULL,
  currency_symbol text NOT NULL,
  UNIQUE (currency_name, currency_symbol)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE currency_names;