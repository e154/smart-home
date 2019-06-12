-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type nodes_status as enum ('enabled', 'disabled');

CREATE TABLE nodes (
  id          bigserial                not null constraint nodes_pkey primary key,
  ip          VarChar(255)             NOT NULL,
  port        numeric                  NOT NULL,
  name        VarChar(255)             NOT NULL,
  description Text                     NOT NULL,
  created_at  timestamp with time zone,
  updated_at  timestamp with time zone null,
  status      nodes_status             NOT NULL DEFAULT 'enabled'
);

CREATE UNIQUE INDEX port_ip_2_nodes_unq
  ON nodes (port, ip);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS nodes CASCADE;
