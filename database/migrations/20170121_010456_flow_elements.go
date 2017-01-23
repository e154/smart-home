package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type FlowElements_20170121_010456 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &FlowElements_20170121_010456{}
	m.Created = "20170121_010456"
	migration.Register("FlowElements_20170121_010456", m)
}

// Run the migrations
func (m *FlowElements_20170121_010456) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE flow_elements (
	uuid VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	flow_id Int( 32 ) NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	prototype_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'default',
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	script_id Int( 11 ) NULL,
	flow_link Int( 32 ) NULL,
	PRIMARY KEY ( uuid ),
	CONSTRAINT id UNIQUE( uuid ),
	CONSTRAINT unique_id UNIQUE( uuid ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *FlowElements_20170121_010456) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `flow_elements` CASCADE")
}
