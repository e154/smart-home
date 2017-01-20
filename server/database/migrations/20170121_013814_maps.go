package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Maps_20170121_013814 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Maps_20170121_013814{}
	m.Created = "20170121_013814"
	migration.Register("Maps_20170121_013814", m)
}

// Run the migrations
func (m *Maps_20170121_013814) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE maps (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	name Char( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	options Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Maps_20170121_013814) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `maps` CASCADE")
}
