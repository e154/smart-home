package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type DropScenarios_20170319_155434 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &DropScenarios_20170319_155434{}
	m.Created = "20170319_155434"
	migration.Register("DropScenarios_20170319_155434", m)
}

// Run the migrations
func (m *DropScenarios_20170319_155434) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scenarios_scenario_scripts`")
	m.SQL("ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scripts_scenario_scripts`")
	m.SQL("DROP TABLE IF EXISTS `scenario_scripts` CASCADE")

	m.SQL("ALTER TABLE `workflows` DROP FOREIGN KEY `lnk_scenarios_workflows`;")
	m.SQL("ALTER TABLE `workflows` DROP COLUMN `scenario_id`")
	m.SQL("DROP TABLE IF EXISTS `scenarios` CASCADE")
}

// Reverse the migrations
func (m *DropScenarios_20170319_155434) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL(`CREATE TABLE scenarios (
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
	ENGINE = InnoDB`)

	m.SQL("INSERT INTO `scenarios` ( `created_at`, `id`, `name`, `system_name`, `update_at`) VALUES ( '2014-06-15 17:50:15', 1, 'default', 'default', '2014-06-15 17:50:15' );")
	m.SQL("ALTER TABLE `workflows` ADD COLUMN `scenario_id` Int( 32 ) NULL;")
	m.SQL("ALTER TABLE `workflows` ADD CONSTRAINT `lnk_scenarios_workflows` FOREIGN KEY ( `scenario_id` ) REFERENCES `scenarios`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")

	m.SQL(`CREATE TABLE scenario_scripts (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	scenario_id Int( 32 ) NOT NULL,
	state VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	script_id Int( 32 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB`)

	m.SQL(`CREATE INDEX lnk_scenarios_scenario_scripts USING BTREE ON scenario_scripts( scenario_id );`)
	m.SQL(`CREATE INDEX lnk_scripts_scenario_scripts USING BTREE ON scenario_scripts( script_id );`)
	m.SQL("ALTER TABLE `scenario_scripts` ADD CONSTRAINT `lnk_scenarios_scenario_scripts` FOREIGN KEY ( `scenario_id` ) REFERENCES `scenarios`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `scenario_scripts` ADD CONSTRAINT `lnk_scripts_scenario_scripts` FOREIGN KEY ( `script_id` ) REFERENCES `scripts`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
}
