-- +migrate Up
create type devices_status as enum ('enabled', 'disabled');
create type devices_type as enum ('default', 'smartbus', 'modbus', 'zigbee');

CREATE TABLE devices (
  id          bigserial constraint devices_pkey primary key not null,
  name        VarChar(255)                                  NOT NULL,
  description VarChar(255)                                  default '',
  device_id   smallint                                      NULL,
  node_id     BIGINT CONSTRAINT devices_2_nodes_fk REFERENCES nodes (id) on update cascade on delete cascade null,
  properties  jsonb                                         default '{"params":{}}',
  type        devices_type                                  default 'default' not null,
  status      devices_status                                NOT NULL DEFAULT 'enabled',
  created_at  timestamp with time zone                      not null,
  updated_at  timestamp with time zone                      not null
);

CREATE UNIQUE INDEX name_device_address_2_devices_unq ON devices (name, device_id);
CREATE INDEX properties_address_2_devices_idx ON devices ((((properties -> 'params') ->> 'address')::int));

CREATE TABLE device_actions (
  id          bigserial constraint device_actions_pkey primary key                                                      NOT NULL,
  device_id   BIGINT CONSTRAINT device_actions_2_devices_fk REFERENCES devices (id) on update cascade on delete cascade,
  name        VarChar(255)                                                                                              NOT NULL,
  description VarChar(255)                                                                                              NULL,
  script_id   BIGINT CONSTRAINT device_actions_2_scripts_fk REFERENCES scripts (id) on update cascade on delete cascade NULL,
  created_at  timestamp with time zone                                                                                  NOT NULL,
  updated_at  timestamp with time zone                                                                                  NOT NULL
);

-- CREATE UNIQUE INDEX device_name_2_device_actions_unq ON device_actions (device_id, name);

CREATE TABLE device_states (
  id          bigserial                not null constraint device_states_pkey primary key,
  system_name VarChar(255)             NOT NULL,
  description Text                     NOT NULL,
  device_id   BIGINT CONSTRAINT device_states_2_devices_fk REFERENCES devices (id) on update cascade on delete cascade,
  created_at  timestamp with time zone not null,
  updated_at  timestamp with time zone not null
);

CREATE UNIQUE INDEX device_nsystem_name_2_device_states_unq ON device_states (device_id, system_name);

-- +migrate Down
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS device_actions CASCADE;
DROP TABLE IF EXISTS device_states CASCADE;
drop type devices_status;
drop type devices_type;