-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- REMEMBER ME TOKENS
CREATE TABLE remember_me_tokens(
  id UUID PRIMARY KEY, 
  selector VARCHAR UNIQUE NOT NULL,
  token_hash VARCHAR NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  valid_to TIMESTAMP NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE remember_me_tokens;