-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table entities
    add column settings jsonb default '{}'::jsonb;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table entities
    drop column settings;

