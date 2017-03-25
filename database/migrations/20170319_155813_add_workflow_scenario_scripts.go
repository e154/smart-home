package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AddWorkflowScenarioScripts_20170319_155813 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddWorkflowScenarioScripts_20170319_155813{}
	m.Created = "20170319_155813"
	migration.Register("AddWorkflowScenarioScripts_20170319_155813", m)
}

// Run the migrations
func (m *AddWorkflowScenarioScripts_20170319_155813) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE workflow_scenario_scripts (
	id Int( 22 ) AUTO_INCREMENT NOT NULL,
	workflow_scenario_id Int( 22 ) NOT NULL,
	script_id Int( 22 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	ENGINE = InnoDB;`)

	m.SQL("ALTER TABLE `workflow_scenario_scripts` ADD CONSTRAINT `lnk_scripts_workflow_scenario_scripts` FOREIGN KEY ( `script_id` ) REFERENCES `scripts`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")

}

// Reverse the migrations
func (m *AddWorkflowScenarioScripts_20170319_155813) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `workflow_scenario_scripts` DROP FOREIGN KEY `lnk_scripts_workflow_scenario_scripts`;")
	m.SQL("DROP TABLE IF EXISTS `workflow_scenario_scripts` CASCADE")
}
