-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type layer_status as enum ('enabled', 'disabled', 'frozen');
create type map_element_status as enum ('enabled', 'disabled', 'frozen');

CREATE TABLE maps (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description Text         NOT NULL,
  options     JSONB DEFAULT '{}',
  created_at  TIMESTAMPTZ  NOT NULL,
  updated_at  TIMESTAMPTZ  NOT NULL
);

CREATE TABLE map_texts (
  id    BIGSERIAL PRIMARY KEY,
  Text  Text NOT NULL,
  style Text DEFAULT '{}'
);

CREATE TABLE map_layers (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description Text         NOT NULL,
  map_id      BIGINT       NULL CONSTRAINT map_layers_2_maps_fk REFERENCES maps (id) ON UPDATE CASCADE ON DELETE RESTRICT,
  status      layer_status NOT NULL,
  weight      INTEGER      NOT NULL DEFAULT 0,
  created_at  TIMESTAMPTZ  NOT NULL,
  updated_at  TIMESTAMPTZ  NOT NULL
);

CREATE TABLE map_images (
  id       BIGSERIAL PRIMARY KEY,
  image_id BIGINT CONSTRAINT map_images_2_images_fk REFERENCES images (id) ON UPDATE CASCADE ON DELETE SET NULL,
  style    Text DEFAULT '{}'
);

CREATE TABLE map_elements (
  id             BIGSERIAL PRIMARY KEY,
  name           VARCHAR(255)       NOT NULL,
  description    Text               NOT NULL,
  prototype_id   Text               NOT NULL,
  prototype_type VARCHAR(255)       NOT NULL,
  graph_settings JSONB              NOT NULL DEFAULT '{}',
  map_layer_id   BIGINT             NULL CONSTRAINT map_elements_2_map_layers_fk REFERENCES map_layers (id) ON UPDATE CASCADE ON DELETE RESTRICT,
  map_id         BIGINT             NULL CONSTRAINT map_elements_2_maps_fk REFERENCES maps (id) ON UPDATE CASCADE ON DELETE RESTRICT,
  status         map_element_status NOT NULL,
  weight         INTEGER            NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ        NOT NULL,
  updated_at     TIMESTAMPTZ        NOT NULL
);

CREATE TABLE map_devices (
  id          BIGSERIAL PRIMARY KEY,
  system_name VARCHAR(255) NOT NULL,
  device_id   BIGINT CONSTRAINT map_devices_2_devices_fk REFERENCES devices (id) ON UPDATE CASCADE ON DELETE SET NULL,
  image_id    BIGINT CONSTRAINT map_devices_2_images_fk REFERENCES images (id) ON UPDATE CASCADE ON DELETE SET NULL,
  created_at  TIMESTAMPTZ  NOT NULL,
  updated_at  TIMESTAMPTZ  NOT NULL
);

CREATE UNIQUE INDEX system_name_device_id_at_map_devices_unq ON map_devices (system_name, device_id);
CREATE INDEX system_name_at_map_devices_idx ON map_devices (system_name);
CREATE INDEX device_at_map_devices_idx ON map_devices (device_id);

CREATE TABLE map_device_states (
  id              BIGSERIAL PRIMARY KEY,
  device_state_id BIGINT CONSTRAINT map_device_states_2_device_states_fk REFERENCES device_states (id) ON UPDATE CASCADE ON DELETE SET NULL,
  map_device_id   BIGINT CONSTRAINT map_device_states_2_map_device_fk REFERENCES map_devices (id) ON UPDATE CASCADE ON DELETE SET NULL,
  image_id        BIGINT CONSTRAINT map_device_states_2_images_fk REFERENCES images (id) ON UPDATE CASCADE ON DELETE SET NULL,
  style           Text DEFAULT '{}',
  created_at      TIMESTAMPTZ NOT NULL,
  updated_at      TIMESTAMPTZ NOT NULL
);

CREATE TABLE map_device_actions (
  id               BIGSERIAL PRIMARY KEY,
  device_action_id BIGINT CONSTRAINT map_device_actions_2_device_actions_fk REFERENCES device_actions (id) ON UPDATE CASCADE ON DELETE SET NULL,
  map_device_id    BIGINT CONSTRAINT map_device_actions_2_map_devices_fk REFERENCES map_devices (id) ON UPDATE CASCADE ON DELETE SET NULL,
  image_id         BIGINT CONSTRAINT map_device_actions_2_images_fk REFERENCES images (id) ON UPDATE CASCADE ON DELETE SET NULL,
  type             Text        NOT NULL,
  created_at       TIMESTAMPTZ NOT NULL,
  updated_at       TIMESTAMPTZ NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS map_texts CASCADE;
DROP TABLE IF EXISTS maps CASCADE;
DROP TABLE IF EXISTS map_layers CASCADE;
DROP TABLE IF EXISTS map_images CASCADE;
DROP TABLE IF EXISTS map_elements CASCADE;
DROP TABLE IF EXISTS map_devices CASCADE;
DROP TABLE IF EXISTS map_device_states CASCADE;
DROP TABLE IF EXISTS map_device_actions CASCADE;
drop type layer_status;
drop type map_element_status;