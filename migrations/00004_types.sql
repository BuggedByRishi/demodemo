-- +goose Up
-- create types for the database
CREATE TYPE core.apps AS ENUM (
    'tallyprime7.0-wine10',
    'tallyprime7.0-wine11',
    'pytha'
);


-- +goose Down
DROP TYPE IF EXISTS core.apps;