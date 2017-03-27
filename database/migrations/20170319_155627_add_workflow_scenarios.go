package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AddWorkflowScenarios_20170319_155627 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddWorkflowScenarios_20170319_155627{}
	m.Created = "20170319_155627"
	migration.Register("AddWorkflowScenarios_20170319_155627", m)
}

// Run the migrations
func (m *AddWorkflowScenarios_20170319_155627) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE workflow_scenarios (
	id Int( 22 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) NOT NULL,
	system_name VarChar( 255 ) NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	workflow_id Int( 22 ) NOT NULL,
	PRIMARY KEY ( id ) )
	ENGINE = InnoDB;`)

	m.SQL("ALTER TABLE `workflow_scenarios` ADD CONSTRAINT `lnk_workflows_workflow_scenarios` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `workflow_scenarios` ADD CONSTRAINT `unique` UNIQUE( `workflow_id`, `system_name` );")
}

// Reverse the migrations
func (m *AddWorkflowScenarios_20170319_155627) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `workflow_scenarios` DROP FOREIGN KEY `lnk_workflows_workflow_scenarios`;")
	m.SQL("DROP TABLE IF EXISTS `workflow_scenarios` CASCADE")
}
