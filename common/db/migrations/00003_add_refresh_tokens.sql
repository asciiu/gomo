-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- REMEMBER ME TOKENS
CREATE TABLE refresh_tokens(
  id UUID PRIMARY KEY, 
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  selector VARCHAR UNIQUE NOT NULL,
  token_hash VARCHAR NOT NULL,
  expires_on TIMESTAMP NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE refresh_tokens;