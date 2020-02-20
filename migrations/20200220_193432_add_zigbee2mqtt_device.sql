-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE zigbee2mqtt_devices
(
    id           text
        constraint zigbee2mqtt_devices_pkey primary key not null,
    name         text,
    type         text,
    model        text,
    manufacturer text,
    functions    text[],
    created_at   timestamp with time zone               not null,
    updated_at   timestamp with time zone               not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table zigbee2mqtt_devices cascade;
