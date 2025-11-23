-- +goose Up
create table user_logs (
    id serial primary key,
    action text not null,
    entity_id int not null,
    created_at timestamp not null default now()
);
-- +goose Down
drop table user_logs;