-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE zigbee2mqtt
(
    id                 bigserial
        constraint zigbee2mqtt_pkey primary key not null,
    name               text,
    login              text                     null,
    encrypted_password text                     null,
    created_at         timestamp with time zone not null,
    updated_at         timestamp with time zone not null
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table zigbee2mqtt cascade;
