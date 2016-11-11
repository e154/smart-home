package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Workflows_20161110_000653 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Workflows_20161110_000653{}
	m.Created = "20161110_000653"
	migration.Register("Workflows_20161110_000653", m)
}

// Run the migrations
func (m *Workflows_20161110_000653) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE workflows (
	id Int( 255 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) NOT NULL,
	status Enum( 'enabled', 'disabled' ) NOT NULL DEFAULT 'enabled',
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	PRIMARY KEY ( id ) )
	ENGINE = InnoDB
	AUTO_INCREMENT = 1;`)
}

// Reverse the migrations
func (m *Workflows_20161110_000653) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `workflows` CASCADE")
}
