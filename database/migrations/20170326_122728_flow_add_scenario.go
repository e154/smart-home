package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type FlowAddScenario_20170326_122728 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &FlowAddScenario_20170326_122728{}
	m.Created = "20170326_122728"
	migration.Register("FlowAddScenario_20170326_122728", m)
}

// Run the migrations
func (m *FlowAddScenario_20170326_122728) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE `flows` ADD COLUMN `workflow_scenario_id` Int( 22 ) NULL;")
}

// Reverse the migrations
func (m *FlowAddScenario_20170326_122728) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `flows` DROP COLUMN `workflow_scenario_id`;")
}
