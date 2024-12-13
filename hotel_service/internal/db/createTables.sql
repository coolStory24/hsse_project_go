create table administrators (
    id uuid primary key,
    nickname text not null,
    phone_number text,
    email text not null,
    password text not null
);

create table hotels (
    id uuid primary key,
    hotel_name text not null,
    night_price int not null,
    administrator_id uuid not null,
    foreign key (administrator_id) references administrators (id)
);
