-- +goose Up
create table chat_logs (
    id serial primary key,
    action text not null,
    entity_id int not null,
    created_at timestamp not null default now()
);
-- +goose Down
drop table chat_logs;