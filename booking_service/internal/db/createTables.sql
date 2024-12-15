create table clients (
    id uuid primary key,
    nickname text not null,
    phone_number text,
    email text not null,
    password text not null
);

create table bookings (
    id uuid primary key default gen_random_uuid(),
    hotel_id uuid not null,
    client_id uuid not null,
    check_in_date time not null,
    check_out_date time not null,
    foreign key (client_id) references clients (id)
);
