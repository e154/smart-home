package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MapDevices_20170121_013341 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MapDevices_20170121_013341{}
	m.Created = "20170121_013341"
	migration.Register("MapDevices_20170121_013341", m)
}

// Run the migrations
func (m *MapDevices_20170121_013341) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE map_devices (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	device_id Int( 32 ) NOT NULL,
	image_id Int( 32 ) NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( device_id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *MapDevices_20170121_013341) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `map_devices` CASCADE")
}
