-- +goose Up
ALTER TABLE accounts ADD COLUMN color text;


-- +goose Down
ALTER TABLE accounts DROP COLUMN color; 
