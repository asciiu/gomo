-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- PASSWORD RESET CODES
CREATE TABLE password_reset_codes(
  id UUID PRIMARY KEY,
  code VARCHAR NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  valid_to TIMESTAMP NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE password_reset_codes;
