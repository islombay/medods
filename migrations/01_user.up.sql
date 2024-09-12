create extension if not exists "uuid-ossp";

create table if not exists users (
    id uuid primary key default uuid_generate_v4(),
    first_name varchar(255),
    last_name varchar(255),
    email varchar(255),

    hash_token varchar(255),

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);

create index index_user_deleted_at on users (deleted_at);