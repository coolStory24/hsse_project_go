-- +goose Up
-- +goose StatementBegin
CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TYPE role_enum AS ENUM ('owner', 'guest');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role role_enum NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role_enum;
-- +goose StatementEnd
