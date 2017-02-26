package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type ScenarioScripts_20170225_233820 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ScenarioScripts_20170225_233820{}
	m.Created = "20170225_233820"
	migration.Register("ScenarioScripts_20170225_233820", m)
}

// Run the migrations
func (m *ScenarioScripts_20170225_233820) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
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

// Reverse the migrations
func (m *ScenarioScripts_20170225_233820) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scenarios_scenario_scripts`;")
	m.SQL("ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scripts_scenario_scripts`;")
	m.SQL("DROP TABLE IF EXISTS `scenario_scripts` CASCADE")
}
