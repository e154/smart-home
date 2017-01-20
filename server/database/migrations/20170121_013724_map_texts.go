package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MapTexts_20170121_013724 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MapTexts_20170121_013724{}
	m.Created = "20170121_013724"
	migration.Register("MapTexts_20170121_013724", m)
}

// Run the migrations
func (m *MapTexts_20170121_013724) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE map_texts (
	id Int( 11 ) AUTO_INCREMENT NOT NULL,
	Text Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	Style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *MapTexts_20170121_013724) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `map_texts` CASCADE")
}
