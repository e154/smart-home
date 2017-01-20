package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Workflows_20170121_005244 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Workflows_20170121_005244{}
	m.Created = "20170121_005244"
	migration.Register("Workflows_20170121_005244", m)
}

// Run the migrations
func (m *Workflows_20170121_005244) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE workflows (
	id Int( 255 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Workflows_20170121_005244) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `workflows` CASCADE")
}
