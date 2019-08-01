-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table nodes
    drop column if exists ip;
alter table nodes
    drop column if exists port;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table nodes
    add column if not exists ip varchar(255);
alter table nodes
    add column if not exists port numeric;

