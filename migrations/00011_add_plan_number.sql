-- +goose Up
-- Keo wants this to be an incremented counter for the user-plan
ALTER TABLE plans ADD COLUMN user_plan_number integer; 

-- +goose Down
ALTER TABLE plans DROP COLUMN user_plan_number;

