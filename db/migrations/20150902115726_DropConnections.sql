-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
drop table if exists connections cascade;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
create table connections (
id serial,
user_id int references users(id) not null,
message_queue text not null,
primary key (id)
);

create index idx_connections_user_id on connections (user_id);
