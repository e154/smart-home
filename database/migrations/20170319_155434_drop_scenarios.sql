-- +migrate Up
ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scenarios_scenario_scripts`;
ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scripts_scenario_scripts`;
DROP TABLE IF EXISTS `scenario_scripts` CASCADE;

ALTER TABLE `workflows` DROP FOREIGN KEY `lnk_scenarios_workflows`;
ALTER TABLE `workflows` DROP COLUMN `scenario_id`;
DROP TABLE IF EXISTS `scenarios` CASCADE;

-- +migrate Down
CREATE TABLE scenarios (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
system_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
created_at DateTime NOT NULL,
update_at DateTime NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ),
CONSTRAINT unique_system_name UNIQUE( system_name ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

INSERT INTO `scenarios` ( `created_at`, `id`, `name`, `system_name`, `update_at`) VALUES ( '2014-06-15 17:50:15', 1, 'default', 'default', '2014-06-15 17:50:15' );
ALTER TABLE `workflows` ADD COLUMN `scenario_id` Int( 32 ) NULL;
ALTER TABLE `workflows` ADD CONSTRAINT `lnk_scenarios_workflows` FOREIGN KEY ( `scenario_id` ) REFERENCES `scenarios`( `id` ) ON DELETE Restrict ON UPDATE Cascade;

CREATE TABLE scenario_scripts (
id Int( 32 ) AUTO_INCREMENT NOT NULL,
scenario_id Int( 32 ) NOT NULL,
state VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
script_id Int( 32 ) NOT NULL,
PRIMARY KEY ( id ),
CONSTRAINT unique_id UNIQUE( id ) )
CHARACTER SET = utf8
COLLATE = utf8_general_ci
ENGINE = InnoDB;

CREATE INDEX lnk_scenarios_scenario_scripts USING BTREE ON scenario_scripts( scenario_id );
CREATE INDEX lnk_scripts_scenario_scripts USING BTREE ON scenario_scripts( script_id );
ALTER TABLE `scenario_scripts` ADD CONSTRAINT `lnk_scenarios_scenario_scripts` FOREIGN KEY ( `scenario_id` ) REFERENCES `scenarios`( `id` ) ON DELETE Cascade ON UPDATE Cascade;
ALTER TABLE `scenario_scripts` ADD CONSTRAINT `lnk_scripts_scenario_scripts` FOREIGN KEY ( `script_id` ) REFERENCES `scripts`( `id` ) ON DELETE Restrict ON UPDATE Cascade;
