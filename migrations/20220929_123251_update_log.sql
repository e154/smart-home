-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table logs
    add column owner text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


