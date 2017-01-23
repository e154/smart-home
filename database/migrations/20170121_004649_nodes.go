package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Nodes_20170121_004649 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Nodes_20170121_004649{}
	m.Created = "20170121_004649"
	migration.Register("Nodes_20170121_004649", m)
}

// Run the migrations
func (m *Nodes_20170121_004649) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE nodes (
	id Int( 255 ) AUTO_INCREMENT NOT NULL,
	ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	port Int( 255 ) NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	status Enum( 'enabled', 'disabled' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'enabled',
	PRIMARY KEY ( id ),
	CONSTRAINT node UNIQUE( port, ip ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Nodes_20170121_004649) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `nodes` CASCADE")
}
