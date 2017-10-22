-- +migrate Up
CREATE TABLE workflow_scenarios (
id Int( 22 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) NOT NULL,
system_name VarChar( 255 ) NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
workflow_id Int( 22 ) NOT NULL,
PRIMARY KEY ( id ) )
ENGINE = InnoDB;

ALTER TABLE `workflow_scenarios` ADD CONSTRAINT `lnk_workflows_workflow_scenarios` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Cascade ON UPDATE Cascade;
ALTER TABLE `workflow_scenarios` ADD CONSTRAINT `unique` UNIQUE( `workflow_id`, `system_name` );

-- +migrate Down
ALTER TABLE `workflow_scenarios` DROP FOREIGN KEY `lnk_workflows_workflow_scenarios`;
DROP TABLE IF EXISTS `workflow_scenarios` CASCADE;