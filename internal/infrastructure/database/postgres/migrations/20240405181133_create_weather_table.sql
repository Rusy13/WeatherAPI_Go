-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS weather (
                                       city_name VARCHAR(255),
    temp FLOAT,
    date DATE,
    data JSONB,
    PRIMARY KEY (city_name, date),
    FOREIGN KEY (city_name) REFERENCES cities (name)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE weather;

-- +goose StatementEnd