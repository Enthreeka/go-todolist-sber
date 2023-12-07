create extension if not exists "uuid-ossp";

create type role as enum ('admin','user');

create table if not exists "user"(
    id uuid DEFAULT uuid_generate_v4(),
    role role not null default 'user',
    login varchar(100) unique not null,
    password varchar(100) unique not null,
    primary key (id)
);

create table if not exists task(
    id int generated always as identity,
    id_user uuid,
    header varchar(150) not null,
    description text not null ,
    created_at timestamp default current_timestamp not null,
    start_date timestamp not null,
    done bool default false,
    primary key (id),
    foreign key (id_user)
            references "user" (id) on delete cascade
);

create table if not exists session(
    token uuid not null,
    user_id uuid,
    expires_at timestamp not null,
    primary key (token),
    foreign key (user_id)
        references "user" (id) on delete cascade
);

update "user" set role = 'admin' where login = '1234567';
