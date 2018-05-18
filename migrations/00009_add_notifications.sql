-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE notifications (
  id UUID PRIMARY KEY NOT NULL,
  notification_type text NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  object_id UUID,
  title text,
  subtitle text,
  description text, 
  timestamp TIMESTAMP DEFAULT now(),
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE notifications;