-- +migrate Up
CREATE TABLE workflow_scripts (
id Int( 22 ) AUTO_INCREMENT NOT NULL,
workflow_id Int( 2 ) NOT NULL,
script_id Int( 22 ) NOT NULL,
weight Int( 10 ) NOT NULL DEFAULT '0',
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
ENGINE = InnoDB;

CREATE UNIQUE INDEX `unique` ON `workflow_scripts`( `workflow_id`, `script_id` );

ALTER TABLE `workflow_scripts` ADD CONSTRAINT `lnk_scripts_workflow_scripts` FOREIGN KEY ( `script_id` ) REFERENCES `scripts`( `id` ) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `workflow_scripts` ADD CONSTRAINT `lnk_workflows_workflow_scripts` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE `workflow_scripts` DROP FOREIGN KEY `lnk_scripts_workflow_scripts`;
ALTER TABLE `workflow_scripts` DROP FOREIGN KEY `lnk_workflows_workflow_scripts`;

DROP TABLE IF EXISTS `workflow_scripts` CASCADE;
