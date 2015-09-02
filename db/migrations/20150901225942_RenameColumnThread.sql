-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table threads rename column with_username to author_username;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table threads rename column author_username to with_username;
