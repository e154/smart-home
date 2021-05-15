-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create table plugins
(
    name    text primary key not null,
    version text             null,
    enabled boolean default null,
    system  boolean default null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table plugins;

