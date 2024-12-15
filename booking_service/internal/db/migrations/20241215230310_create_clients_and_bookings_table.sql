-- +goose Up
-- +goose StatementBegin


CREATE TABLE clients (
    id UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    phone_number TEXT,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL,
    client_id UUID NOT NULL,
    CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES clients (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
