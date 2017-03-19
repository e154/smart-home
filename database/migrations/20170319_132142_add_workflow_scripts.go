package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AddWorkflowScripts_20170319_132142 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddWorkflowScripts_20170319_132142{}
	m.Created = "20170319_132142"
	migration.Register("AddWorkflowScripts_20170319_132142", m)
}

// Run the migrations
func (m *AddWorkflowScripts_20170319_132142) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`
	CREATE TABLE workflow_scripts (
	id Int( 22 ) NOT NULL,
	workflow_id Int( 2 ) NOT NULL,
	script_id Int( 22 ) NOT NULL,
	weight Int( 10 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	ENGINE = InnoDB;`)

	m.SQL("CREATE UNIQUE INDEX `unique` ON `workflow_scripts`( `workflow_id`, `script_id` );")

	m.SQL("ALTER TABLE `workflow_scripts` ADD CONSTRAINT `lnk_scripts_workflow_scripts` FOREIGN KEY ( `script_id` ) REFERENCES `scripts`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
	m.SQL("ALTER TABLE `workflow_scripts` ADD CONSTRAINT `lnk_workflows_workflow_scripts` FOREIGN KEY ( `workflow_id` ) REFERENCES `workflows`( `id` ) ON DELETE Restrict ON UPDATE Cascade;")
}

// Reverse the migrations
func (m *AddWorkflowScripts_20170319_132142) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `workflow_scripts` DROP FOREIGN KEY `lnk_scripts_workflow_scripts`;")
	m.SQL("ALTER TABLE `workflow_scripts` DROP FOREIGN KEY `lnk_workflows_workflow_scripts`;")

	m.SQL("DROP TABLE IF EXISTS `workflow_scripts` CASCADE")
}
