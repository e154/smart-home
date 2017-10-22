-- +migrate Up
ALTER TABLE `workflows` ADD COLUMN `workflow_scenario_id` Int( 22 ) NULL;
ALTER TABLE `workflows` ADD CONSTRAINT `lnk_workflow_scenarios_workflows` FOREIGN KEY ( `workflow_scenario_id` ) REFERENCES `workflow_scenarios`( `id` ) ON DELETE Cascade ON UPDATE Cascade;

-- +migrate Down
ALTER TABLE `workflows` DROP FOREIGN KEY `lnk_workflow_scenarios_workflows`;
ALTER TABLE `workflows` DROP COLUMN `workflow_scenario_id`;

