create extension if not exists "uuid-ossp";

create table if not exists "user"(
    id uuid DEFAULT uuid_generate_v4(),
    login varchar(100) unique,
    primary key (id)
);

create table if not exists task(
    id int generated always as identity,
    id_user uuid,
    header varchar(150) not null,
    description text not null ,
    created_at timestamp not null,
    start_date timestamp not null,
    done bool default false,
    primary key (id),
    foreign key (id_user)
            references "user" (id) on delete cascade
);

