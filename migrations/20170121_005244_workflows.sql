-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
create type scripts_lang as enum ('ts', 'coffeescript', 'javascript');

CREATE TABLE scripts (
  id          bigserial not null constraint scripts_pkey primary key,
  lang        scripts_lang NOT NULL DEFAULT 'javascript',
  name        VarChar(255) NOT NULL,
  source      Text NULL,
  description Text NULL,
  compiled    Text NULL,
  created_at  timestamp with time zone not null,
  updated_at   timestamp with time zone not null
);

create type workflows_status as enum ('enabled', 'disabled');

CREATE TABLE workflows (
  id                   bigserial                not null constraint workflows_pkey primary key,
  name                 VarChar(255)             NOT NULL,
  status               workflows_status         NOT NULL DEFAULT 'enabled',
  description          Text                     NULL,
  created_at           timestamp with time zone not null,
  updated_at            timestamp with time zone not null
);

CREATE TABLE workflow_scenarios (
  id          bigserial                not null constraint workflow_scenarios_pkey primary key,
  workflow_id BIGINT CONSTRAINT workflow_scenarios_2_workflows_fk REFERENCES workflows (id) on update cascade on delete cascade,
  name        VarChar(255)             NOT NULL,
  system_name VarChar(255)             NOT NULL,
  created_at  timestamp with time zone not null,
  updated_at   timestamp with time zone not null
);

CREATE UNIQUE INDEX system_name_2_workflow_scenarios_unq
  ON workflow_scenarios (system_name);

CREATE TABLE workflow_scenario_scripts (
  id                   bigserial not null constraint workflow_scenario_scripts_pkey primary key,
  workflow_scenario_id BIGINT CONSTRAINT workflow_scenarios_2_workflows_fk REFERENCES workflow_scenarios (id) on update cascade on delete restrict,
  script_id            BIGINT CONSTRAINT script_scenarios_2_scripts_fk REFERENCES scripts (id) on update cascade on delete cascade
);

--
CREATE TABLE workflow_scripts (
  id          bigserial not null constraint workflow_scripts_pkey primary key,
  workflow_id BIGINT CONSTRAINT workflow_scenarios_2_workflows_fk REFERENCES workflows (id) on update cascade on delete cascade,
  script_id   BIGINT CONSTRAINT script_scenarios_2_scripts_fk REFERENCES scripts (id) on update cascade on delete cascade,
  weight      int       NOT NULL
);

CREATE UNIQUE INDEX name_2_workflows_unq
  ON workflows (name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS workflows CASCADE;
DROP TABLE IF EXISTS workflow_scripts CASCADE;
DROP TABLE IF EXISTS workflow_scenarios CASCADE;
DROP TABLE IF EXISTS workflow_scenario_scripts CASCADE;
DROP TABLE IF EXISTS scripts CASCADE;
drop type workflows_status;
drop type scripts_lang;