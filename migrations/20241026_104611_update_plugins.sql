-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table plugins
    add column external bool not null default false;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table plugins
    drop column external;
