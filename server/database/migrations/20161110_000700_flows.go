package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Flows_20161110_000700 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Flows_20161110_000700{}
	m.Created = "20161110_000700"
	migration.Register("Flows_20161110_000700", m)
}

// Run the migrations
func (m *Flows_20161110_000700) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
m.SQL(`CREATE TABLE flows (
	id Int( 32 ) NOT NULL,
	name VarChar( 255 ) NOT NULL,
	status Enum( 'enabled', 'disabled' ) NOT NULL DEFAULT 'enabled',
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	workflow_id Int( 11 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	ENGINE = InnoDB
	AUTO_INCREMENT = 1;`)

m.SQL(`ALTER TABLE flows
	ADD CONSTRAINT lnk_flows_workflows FOREIGN KEY ( workflow_id )
	REFERENCES workflows( id )
	ON DELETE Restrict
	ON UPDATE Restrict;`)
}

// Reverse the migrations
func (m *Flows_20161110_000700) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `flows` DROP FOREIGN KEY `lnk_flows_workflows`;")
	m.SQL("DROP TABLE IF EXISTS `flows` CASCADE")
}
