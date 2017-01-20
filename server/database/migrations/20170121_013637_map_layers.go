package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MapLayers_20170121_013637 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MapLayers_20170121_013637{}
	m.Created = "20170121_013637"
	migration.Register("MapLayers_20170121_013637", m)
}

// Run the migrations
func (m *MapLayers_20170121_013637) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE map_layers (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	map_id Int( 11 ) NOT NULL,
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	weight Int( 11 ) NOT NULL DEFAULT '0',
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *MapLayers_20170121_013637) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `map_layers` CASCADE")
}
