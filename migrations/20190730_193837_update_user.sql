-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table users
    drop column if exists authentication_token;
alter table users
    add column authentication_token text default null;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


