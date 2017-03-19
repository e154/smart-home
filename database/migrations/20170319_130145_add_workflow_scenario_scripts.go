package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AddWorkflowScenarioScripts_20170319_130145 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddWorkflowScenarioScripts_20170319_130145{}
	m.Created = "20170319_130145"
	migration.Register("AddWorkflowScenarioScripts_20170319_130145", m)
}

// Run the migrations
func (m *AddWorkflowScenarioScripts_20170319_130145) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE workflow_scenario_scripts (
	id Int( 22 ) AUTO_INCREMENT NOT NULL,
	workflow_id Int( 22 ) NOT NULL,
	scenario_id Int( 22 ) NOT NULL,
	script_id Int( 22 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	ENGINE = InnoDB;`)

	m.SQL("CREATE UNIQUE INDEX `unique` ON `workflow_scenario_scripts`( `workflow_id`, `scenario_id`, `script_id` );")
	m.SQL("ALTER TABLE `workflow_scenario_scripts` ADD CONSTRAINT `lnk_scenarios_workflow_scenario_scripts` FOREIGN KEY ( `scenario_id` ) REFERENCES `scenarios`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `workflow_scenario_scripts` ADD CONSTRAINT `lnk_scripts_workflow_scenario_scripts` FOREIGN KEY ( `script_id` ) REFERENCES `scripts`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `workflow_scenario_scripts` ADD CONSTRAINT `lnk_workflows_workflow_scenario_scripts` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")

	m.SQL("ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scenarios_scenario_scripts`;")
	m.SQL("ALTER TABLE `scenario_scripts` DROP FOREIGN KEY `lnk_scripts_scenario_scripts`;")
	m.SQL("DROP TABLE IF EXISTS `scenario_scripts` CASCADE")
}

// Reverse the migrations
func (m *AddWorkflowScenarioScripts_20170319_130145) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

	m.SQL("ALTER TABLE `workflow_scenario_scripts` DROP FOREIGN KEY `lnk_scenarios_workflow_scenario_scripts`;")
	m.SQL("ALTER TABLE `workflow_scenario_scripts` DROP FOREIGN KEY `lnk_scripts_workflow_scenario_scripts`;")
	m.SQL("ALTER TABLE `workflow_scenario_scripts` DROP FOREIGN KEY `lnk_workflows_workflow_scenario_scripts`;")
	m.SQL("DROP TABLE IF EXISTS `workflow_scenario_scripts` CASCADE")

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
