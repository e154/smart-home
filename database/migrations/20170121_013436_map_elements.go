package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MapElements_20170121_013436 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MapElements_20170121_013436{}
	m.Created = "20170121_013436"
	migration.Register("MapElements_20170121_013436", m)
}

// Run the migrations
func (m *MapElements_20170121_013436) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE map_elements (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	prototype_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	layer_id Int( 32 ) NOT NULL,
	map_id Int( 32 ) NOT NULL,
	graph_settings Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	weight Int( 11 ) NOT NULL DEFAULT '0',
	prototype_id Int( 32 ) NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *MapElements_20170121_013436) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `map_elements` CASCADE")
}
