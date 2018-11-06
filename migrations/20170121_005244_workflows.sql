-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type workflows_status as enum ('enabled', 'disabled');

CREATE TABLE workflows (
  id          bigserial                not null constraint workflows_pkey primary key,
  name        VarChar(255)             NOT NULL,
  status      workflows_status         NOT NULL DEFAULT 'enabled',
  created_at  timestamp with time zone not null,
  update_at   timestamp with time zone not null,
  description Text                     NULL
);

CREATE UNIQUE INDEX name_2_workflows_unq
  ON workflows (name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS workflows CASCADE;
