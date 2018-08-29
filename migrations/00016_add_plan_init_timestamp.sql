-- +goose Up
-- this shall be the currency used to measure initial and final currency valuations 
ALTER TABLE plans ADD COLUMN initial_timestamp text; 

-- +goose Down
ALTER TABLE plans DROP COLUMN initial_timestamp;