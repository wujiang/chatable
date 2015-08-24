-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table users (
       id serial,
       username text not null,
       first_name text,
       last_name text,
       email text not null,
       phone_number text not null,
       password text not null,
       is_active boolean not null default true,
       created_at timestamp without time zone default now(),
       deactivated_at timestamp without time zone,
       original_ip text,
       user_class text,
       primary key (id),
       unique (username),
       unique (email),
       unique (phone_number)
);


create table threads (
       id serial,
       user_id int references users(id) not null,
       with_user_id int references users(id) not null,
       with_username text not null,
       created_at timestamp without time zone default now(),
       latest_message text not null,
       primary key (id),
       unique (user_id, with_user_id)
);

create index idx_inboxes_user_id_created_at on inboxes (user_id, created_at desc);


create table envelopes (
       id serial,
       user_id int references users(id) not null,
       with_user_id int references users(id) not null,
       is_incoming boolean not null,
       created_at timestamp without time zone default now(),
       deleted_at timestamp without time zone,
       read_at timestamp without time zone,
       message text not null,
       message_type int not null default 0,
       primary key (id)
);

create index idx_envelopes_user_id_with_user_id on envelopes (user_id, with_user_id, deleted_at, created_at desc);


create table connections (
       id serial,
       user_id int references users(id) not null,
       message_queue text not null,
       primary key (id)
);

create index idx_connections_user_id on connections (user_id);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table envelopes;
drop table connections;
drop table threads;
drop table users;
