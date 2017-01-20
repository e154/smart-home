package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Scripts_20170121_014057 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Scripts_20170121_014057{}
	m.Created = "20170121_014057"
	migration.Register("Scripts_20170121_014057", m)
}

// Run the migrations
func (m *Scripts_20170121_014057) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE scripts (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	lang Enum( 'lua', 'coffeescript', 'javascript' ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'lua',
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	source Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	description Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	compiled Text CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ),
	CONSTRAINT unique_name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Scripts_20170121_014057) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `scripts` CASCADE")
}
