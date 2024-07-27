-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table plugins
    add column triggers bool not null default false;

update plugins
set triggers = true
where name in ('time', 'state_change', 'system', 'alexa', 'ble');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table plugins
    drop column triggers;

