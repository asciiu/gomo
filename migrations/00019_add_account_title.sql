-- +goose Up
ALTER TABLE accounts ADD COLUMN title text;


-- +goose Down
ALTER TABLE accounts DROP COLUMN title; 
