-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table map_devices
    drop column system_name cascade;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table map_devices
    drop column system_name cascade;

CREATE INDEX system_name_at_map_devices_idx ON map_devices (system_name);

