create table hotels (
    id uuid primary key,
    hotel_name varchar(200) not null,
    night_price int not null,
    administrator_id uuid not null,
    foreign key (administrator_id) references administrators (id)
);

create table administrators (
    id uuid primary key,
    nickname varchar(50) not null,
    phone_number varchar(20),
    email varchar(100) not null,
    password varchar(150) not null
);
