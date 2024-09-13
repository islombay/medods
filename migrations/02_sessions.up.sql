create table if not exists sessions(
    id uuid primary key default uuid_generate_v4() not null,
    hash varchar(255),
    user_id uuid,
    ip varchar(20),

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (user_id) references users(id)
);

create index index_session_hash on sessions (hash);
create index index_session_deleted_at on sessions (deleted_at);
