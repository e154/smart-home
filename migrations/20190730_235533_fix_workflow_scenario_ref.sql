-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE workflows
    DROP CONSTRAINT if exists workflow_2_workflow_scenarios_fk;
ALTER TABLE workflows
    ADD CONSTRAINT workflow_2_workflow_scenarios_fk FOREIGN KEY (workflow_scenario_id)
        REFERENCES workflow_scenarios (id) ON UPDATE CASCADE ON DELETE SET NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

