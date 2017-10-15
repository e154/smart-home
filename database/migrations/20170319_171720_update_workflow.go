package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UpdateWorkflow_20170319_171720 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UpdateWorkflow_20170319_171720{}
	m.Created = "20170319_171720"
	migration.Register("UpdateWorkflow_20170319_171720", m)
}

// Run the migrations
func (m *UpdateWorkflow_20170319_171720) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE `workflows` ADD COLUMN `workflow_scenario_id` Int( 22 ) NULL;")
	m.SQL("ALTER TABLE `workflows` ADD CONSTRAINT `lnk_workflow_scenarios_workflows` FOREIGN KEY ( `workflow_scenario_id` ) REFERENCES `workflow_scenarios`( `id` ) ON DELETE Cascade ON UPDATE Cascade;")
}

// Reverse the migrations
func (m *UpdateWorkflow_20170319_171720) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `workflows` DROP FOREIGN KEY `lnk_workflow_scenarios_workflows`;")
	m.SQL("ALTER TABLE `workflows` DROP COLUMN `workflow_scenario_id`;")
}
