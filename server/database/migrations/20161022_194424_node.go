package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Node_20161022_194424 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Node_20161022_194424{}
	m.Created = "20161022_194424"
	migration.Register("Node_20161022_194424", m)
}

// Run the migrations
func (m *Node_20161022_194424) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE node (
	id Int( 255 ) AUTO_INCREMENT NOT NULL,
	ip VarChar( 255 ) NOT NULL,
	port Int( 255 ) NOT NULL,
	name VarChar( 255 ) NOT NULL,
	description Text NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	status Enum( 'enabled', 'disabled' ) NOT NULL DEFAULT 'enabled',
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ),
	CONSTRAINT node UNIQUE( port, ip ) )
	ENGINE = InnoDB
	AUTO_INCREMENT = 1;`)
}

// Reverse the migrations
func (m *Node_20161022_194424) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `node` CASCADE")
}
