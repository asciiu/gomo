-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE user_devices (
 id UUID PRIMARY KEY NOT NULL,
 user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
 device_id VARCHAR NOT NULL,
 device_type VARCHAR NOT NULL,
 device_token VARCHAR NOT NULL,
 created_on TIMESTAMP DEFAULT now(),
 updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user_devices;