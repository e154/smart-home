package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type DeviceActions_20170121_012048 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &DeviceActions_20170121_012048{}
	m.Created = "20170121_012048"
	migration.Register("DeviceActions_20170121_012048", m)
}

// Run the migrations
func (m *DeviceActions_20170121_012048) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE device_actions (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	device_id Int( 11 ) NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	script_id Int( 32 ) NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *DeviceActions_20170121_012048) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `device_actions` CASCADE")
}
