package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Logs_20170121_012936 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Logs_20170121_012936{}
	m.Created = "20170121_012936"
	migration.Register("Logs_20170121_012936", m)
}

// Run the migrations
func (m *Logs_20170121_012936) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE logs (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	body Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	level Enum( 'Emergency', 'Alert', 'Critical', 'Error', 'Warning', 'Notice', 'Info', 'Debug' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'Info',
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Logs_20170121_012936) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `logs` CASCADE")
}
