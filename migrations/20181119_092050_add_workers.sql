-- +migrate Up
create type workers_status as enum ('enabled', 'disabled');

CREATE TABLE workers (
  id               bigserial                not null constraint workers_pkey primary key,
  name             VarChar(255)             NOT NULL,
  time             VarChar(254)             NOT NULL,
  status           devices_status           NOT NULL DEFAULT 'enabled',
  flow_id          BIGINT CONSTRAINT workers_2_flows_fk REFERENCES flows (id) on update cascade on delete cascade,
  workflow_id      BIGINT CONSTRAINT workers_2_workflows_fk REFERENCES workflows (id) on update cascade on delete cascade,
  device_action_id BIGINT CONSTRAINT workers_2_device_actions_fk REFERENCES device_actions (id) on update cascade on delete cascade,
  created_at       timestamp with time zone not null,
  updated_at       timestamp with time zone not null
);

-- +migrate Down
DROP TABLE IF EXISTS workers CASCADE;
drop type workers_status;
