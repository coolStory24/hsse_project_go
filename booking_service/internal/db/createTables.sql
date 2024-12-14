create table clients (
    id uuid primary key,
    nickname text not null,
    phone_number text,
    email text not null,
    password text not null
);

create table bookings (
    id uuid primary key,
    room_id uuid not null,
    client_id uuid not null,
    foreign key (client_id) references clients (id)
);
