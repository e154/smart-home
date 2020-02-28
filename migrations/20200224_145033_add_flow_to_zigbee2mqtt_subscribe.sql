-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

create table flow_zigbee2mqtt_devices
(
    id                    BIGSERIAL PRIMARY KEY,
    flow_id               BIGINT      NOT NULL
        CONSTRAINT flow_zigbee2mqtt_devices_2_flows_fk REFERENCES flows (id) ON UPDATE CASCADE ON DELETE CASCADE,
    zigbee2mqtt_device_id text        NOT NULL
        CONSTRAINT flow_zigbee2mqtt_devices_2_zigbee2mqtt_device_fk REFERENCES zigbee2mqtt_devices (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at            TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX topic_at_flow_zigbee2mqtt_devices_unq
    ON flow_zigbee2mqtt_devices (flow_id, zigbee2mqtt_device_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table if exists flow_zigbee2mqtt_devices cascade;

