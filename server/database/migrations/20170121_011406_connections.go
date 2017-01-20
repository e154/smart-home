package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Connections_20170121_011406 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Connections_20170121_011406{}
	m.Created = "20170121_011406"
	migration.Register("Connections_20170121_011406", m)
}

// Run the migrations
func (m *Connections_20170121_011406) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE connections (
	uuid VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	element_from VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	element_to VarChar( 254 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	flow_id Int( 32 ) NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	point_from Int( 11 ) NOT NULL,
	point_to Int( 11 ) NOT NULL,
	direction VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( uuid ),
	CONSTRAINT unique_id UNIQUE( uuid ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Connections_20170121_011406) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `connections` CASCADE")
}
