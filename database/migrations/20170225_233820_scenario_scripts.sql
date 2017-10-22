-- +migrate Up
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

-- +migrate Down
ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scenarios_scenario_scripts`;
ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scripts_scenario_scripts`;
DROP TABLE IF EXISTS `scenario_scripts` CASCADE;
