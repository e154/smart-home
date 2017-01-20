package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Images_20170121_012702 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Images_20170121_012702{}
	m.Created = "20170121_012702"
	migration.Register("Images_20170121_012702", m)
}

// Run the migrations
func (m *Images_20170121_012702) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE images (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	thumb VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	image VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	mime_type VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	title VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	size Int( 11 ) NOT NULL,
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	created_at DateTime NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Images_20170121_012702) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `images` CASCADE")
}
