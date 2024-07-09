-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
id SERIAL PRIMARY KEY,
username VARCHAR(255) NOT NULL UNIQUE,
password VARCHAR(255) NOT NULL
);

CREATE TABLE favorite_cities (
 user_id INT REFERENCES users(id),
 city_name VARCHAR(255) NOT NULL,
 PRIMARY KEY (user_id, city_name)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users;
DROP TABLE favorite_cities;

-- +goose StatementEnd