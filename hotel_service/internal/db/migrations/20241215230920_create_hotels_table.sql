-- +goose Up
-- +goose StatementBegin
CREATE TABLE administrators (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    phone_number TEXT,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE hotels (
    id UUID PRIMARY KEY,
    hotel_name TEXT NOT NULL,
    night_price INT NOT NULL,
    administrator_id UUID NOT NULL,
    CONSTRAINT fk_administrator FOREIGN KEY (administrator_id) REFERENCES administrators (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hotels;
DROP TABLE IF EXISTS administrators;
-- +goose StatementEnd
