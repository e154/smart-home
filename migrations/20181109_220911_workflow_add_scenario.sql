-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE public.workflows
  ADD COLUMN workflow_scenario_id BIGINT CONSTRAINT workflow_2_workflow_scenarios_fk REFERENCES workflow_scenarios (id) on update cascade on delete restrict;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table public.workflows
  drop column workflow_scenario_id cascade;

