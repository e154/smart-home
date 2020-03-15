-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE map_device_history
(
    id            BIGSERIAL PRIMARY KEY,
    map_device_id BIGINT                   NOT NULL
        CONSTRAINT map_device_history_2_map_devices_fk REFERENCES map_devices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    type          level_type               NOT NULL DEFAULT 'Info',
    description   text                     NOT NULL,
    created_at    timestamp with time zone not null
);

create index type_at_map_device_history_idx on map_device_history (type);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table map_device_history cascade;

