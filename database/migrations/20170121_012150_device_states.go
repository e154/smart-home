package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type DeviceStates_20170121_012150 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &DeviceStates_20170121_012150{}
	m.Created = "20170121_012150"
	migration.Register("DeviceStates_20170121_012150", m)
}

// Run the migrations
func (m *DeviceStates_20170121_012150) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE device_states (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	system_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	device_id Int( 11 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_1 UNIQUE( device_id, system_name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *DeviceStates_20170121_012150) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `device_actions` CASCADE")
}
