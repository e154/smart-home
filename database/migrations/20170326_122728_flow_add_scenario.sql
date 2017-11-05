-- +migrate Up
ALTER TABLE `flows` ADD COLUMN `workflow_scenario_id` Int( 22 ) NULL;
ALTER TABLE `flows` ADD CONSTRAINT `lnk_flows_workflow_scenarios` FOREIGN KEY (`workflow_scenario_id`) REFERENCES `workflow_scenarios` (`id`);

-- +migrate Down
ALTER TABLE `flows` DROP FOREIGN KEY `lnk_flows_workflow_scenarios`;
ALTER TABLE `flows` DROP COLUMN `workflow_scenario_id`;
