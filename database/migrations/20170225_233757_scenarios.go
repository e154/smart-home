package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Scenarios_20170225_233757 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Scenarios_20170225_233757{}
	m.Created = "20170225_233757"
	migration.Register("Scenarios_20170225_233757", m)
}

// Run the migrations
func (m *Scenarios_20170225_233757) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
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
}

// Reverse the migrations
func (m *Scenarios_20170225_233757) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `workflows` DROP FOREIGN KEY `lnk_scenarios_workflows`;")
	m.SQL("ALTER TABLE `workflows` DROP COLUMN `scenario_id`;")
	m.SQL("DROP TABLE IF EXISTS `scenarios` CASCADE")
}
