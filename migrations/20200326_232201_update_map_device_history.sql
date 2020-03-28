-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table map_device_history
    rename column type to log_level;

alter table map_device_history
    add column type text;

update map_device_history
set type = 'state'
where type = '' or type isnull ;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table map_device_history
    drop column type cascade;

alter table map_device_history
    rename column log_level to type;
