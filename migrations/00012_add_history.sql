-- +goose Up
DROP TABLE notifications;

CREATE TABLE history (
  id UUID PRIMARY KEY NOT NULL,
  "type" text NOT NULL,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  object_id UUID,
  title text,
  subtitle text,
  description text, 
  click_at TIMESTAMP,           -- nullable timestamps of when user clicked and saw
  seen_at TIMESTAMP,            
  timestamp TIMESTAMP DEFAULT now(),
  created_on TIMESTAMP DEFAULT now(),
  updated_on TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
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

DROP TABLE history;

