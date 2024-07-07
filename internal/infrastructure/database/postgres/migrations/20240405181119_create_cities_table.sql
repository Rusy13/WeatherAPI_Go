-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cities (
    name VARCHAR(255) PRIMARY KEY,
    country VARCHAR(255),
    latitude FLOAT,
    longitude FLOAT
    );



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cities;
-- +goose StatementEnd
