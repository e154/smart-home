package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Dashboards_20170121_011938 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Dashboards_20170121_011938{}
	m.Created = "20170121_011938"
	migration.Register("Dashboards_20170121_011938", m)
}

// Run the migrations
func (m *Dashboards_20170121_011938) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE dashboards (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	widgets Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Dashboards_20170121_011938) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `dashboards` CASCADE")
}
