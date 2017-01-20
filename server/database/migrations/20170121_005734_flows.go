package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Flows_20170121_005734 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Flows_20170121_005734{}
	m.Created = "20170121_005734"
	migration.Register("Flows_20170121_005734", m)
}

// Run the migrations
func (m *Flows_20170121_005734) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE flows (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	workflow_id Int( 11 ) NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Flows_20170121_005734) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `flows` CASCADE")
}
