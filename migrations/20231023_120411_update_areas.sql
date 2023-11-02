-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table areas
    drop column if exists polygon cascade,
    add column polygon polygon;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

