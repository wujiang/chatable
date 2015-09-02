-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table auth_tokens add constraint auth_tokens_user_id_fkey foreign key (user_id) references users (id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table auth_tokens drop constraint auth_tokens_user_id_fkey;
