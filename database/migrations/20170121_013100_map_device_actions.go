package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MapDeviceActions_20170121_013100 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MapDeviceActions_20170121_013100{}
	m.Created = "20170121_013100"
	migration.Register("MapDeviceActions_20170121_013100", m)
}

// Run the migrations
func (m *MapDeviceActions_20170121_013100) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE map_device_actions (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	device_action_id Int( 32 ) NOT NULL,
	map_device_id Int( 2 ) NOT NULL,
	image_id Int( 32 ) NULL,
	type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *MapDeviceActions_20170121_013100) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `map_device_actions` CASCADE")
}
