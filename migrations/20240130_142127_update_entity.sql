-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table entities
    add column restore_state bool NOT NULL DEFAULT true;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


