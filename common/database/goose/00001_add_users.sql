-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- USERS
CREATE TABLE users(
    id uuid PRIMARY KEY,
    first_name VARCHAR,
    last_name VARCHAR,
    email VARCHAR UNIQUE NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    password_hash VARCHAR NOT NULL,
    salt VARCHAR NOT NULL,
    created_on TIMESTAMP DEFAULT now(),
    updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;