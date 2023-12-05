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
    created_at timestamp default current_timestamp not null,
    start_date timestamp not null,
    done bool default false,
    primary key (id),
    foreign key (id_user)
            references "user" (id) on delete cascade
);

update task set header = 'Second header' ,description = 'First description' where id = 5 returning *;

select * from task where id_user = '53153c2c-1c10-4b92-b5ff-0cf67b116654';

select * from task
where id_user = '53153c2c-1c10-4b92-b5ff-0cf67b116654'
order by id desc
LIMIT 3;
-- offset 3;

-- 0       3       6
-- 10 9 8| 7 6 5 | 2

-- 1 - 0, 2 - 3, 3 - 6, 4 - 9, 5 - 12, 6 - 15

-- offset = (page-1) * 3

-- 2 страницу - 3 | 3 страницу - 6 | 4 страницу - 9
select * from task
where id_user = '53153c2c-1c10-4b92-b5ff-0cf67b116654' and done = false
order by id desc
offset 0
limit 3;

select * from task
where id_user = '53153c2c-1c10-4b92-b5ff-0cf67b116654' and id < 5
order by id desc
limit 3;