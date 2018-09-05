-- +goose Up
ALTER TABLE accounts ADD COLUMN account_type text;


-- +goose Down
ALTER TABLE accounts DROP COLUMN account_type; 
