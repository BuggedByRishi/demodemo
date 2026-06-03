-- +goose Up
-- case insensitive character string type
-- used to store emails, user attributes, template language codes, etc
CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA extensions;


-- +goose Down
DROP EXTENSION IF EXISTS citext;