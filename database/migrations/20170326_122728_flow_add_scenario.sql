-- +migrate Up
ALTER TABLE `flows` ADD COLUMN `workflow_scenario_id` Int( 22 ) NULL;

-- +migrate Down
ALTER TABLE `flows` DROP COLUMN `workflow_scenario_id`;
