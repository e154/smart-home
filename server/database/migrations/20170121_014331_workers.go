package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Workers_20170121_014331 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Workers_20170121_014331{}
	m.Created = "20170121_014331"
	migration.Register("Workers_20170121_014331", m)
}

// Run the migrations
func (m *Workers_20170121_014331) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE workers (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	flow_id Int( 11 ) NOT NULL,
	workflow_id Int( 11 ) NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	device_action_id Int( 11 ) NOT NULL,
	time VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	PRIMARY KEY ( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Workers_20170121_014331) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `workers` CASCADE")
}
