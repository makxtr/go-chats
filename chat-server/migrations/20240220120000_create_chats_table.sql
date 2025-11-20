-- +goose Up
create table chats (
    id serial primary key,
    usernames text[] not null,
    created_at timestamp not null default now()
);

-- +goose Down
drop table chats;