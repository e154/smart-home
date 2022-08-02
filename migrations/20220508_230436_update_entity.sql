-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table entities
    rename column parent to parent_id;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


