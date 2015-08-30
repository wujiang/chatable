-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table auth_tokens (
       id serial,
       access_key_id text not null unique,
       secret_access_key text not null,
       refresh_token text not null,
       created_at timestamp without time zone default now(),
       expires_at timestamp without time zone,
       modified_at timestamp without time zone default now(),
       user_id int,
       client_id int,
       is_active boolean,
       is_refreshable boolean,
       scope text[],
       primary key (id),
       unique (access_key_id)
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table auth_tokens;
