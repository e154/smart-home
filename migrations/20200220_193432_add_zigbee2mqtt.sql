-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE zigbee2mqtt
(
    id                 bigserial
        constraint zigbee2mqtt_pkey primary key not null,
    name               text,
    login              text                     null,
    encrypted_password text                     null,
    permit_join        boolean default true,
    base_topic         text    default 'zigbee2mqtt',
    created_at         timestamp with time zone not null,
    updated_at         timestamp with time zone not null
);

create unique index base_topic_at_zigbee2mqtt_unq on zigbee2mqtt (base_topic);
create type zigbee2mqtt_devices_status as enum ('active', 'banned', 'removed');
create index login_at_zigbee2mqtt_idx on zigbee2mqtt (login);

CREATE TABLE zigbee2mqtt_devices
(
    id             text
        constraint zigbee2mqtt_devices_pkey primary key not null,
    zigbee2mqtt_id bigint
        CONSTRAINT zigbee2mqtt_devices_2_zigbee2mqtt_fk REFERENCES zigbee2mqtt (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    name           text,
    type           text,
    model          text,
    description    text,
    manufacturer   text,
    functions      text[],
    status         zigbee2mqtt_devices_status default 'active',
    created_at     timestamp with time zone             not null,
    updated_at     timestamp with time zone             not null
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table zigbee2mqtt cascade;
drop table zigbee2mqtt_devices cascade;
drop type zigbee2mqtt_devices_status cascade;

