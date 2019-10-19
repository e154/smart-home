-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE workflow_scenario_scripts
    RENAME CONSTRAINT workflow_scenarios_2_workflows_fk TO workflow_scenario_scripts_2_workflows_fk;
ALTER TABLE workflow_scenario_scripts
    RENAME CONSTRAINT script_scenarios_2_scripts_fk TO workflow_scenario_scripts_2_scripts_fk;

ALTER TABLE workflow_scripts
    RENAME CONSTRAINT workflow_scenarios_2_workflows_fk TO workflow_scripts_2_workflows_fk;
ALTER TABLE workflow_scripts
    RENAME CONSTRAINT script_scenarios_2_scripts_fk TO workflow_scripts_2_scripts_fk;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back


