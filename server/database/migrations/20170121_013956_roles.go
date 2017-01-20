package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Roles_20170121_013956 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Roles_20170121_013956{}
	m.Created = "20170121_013956"
	migration.Register("Roles_20170121_013956", m)
}

// Run the migrations
func (m *Roles_20170121_013956) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE roles (
	name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	parent VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	PRIMARY KEY ( name ),
	CONSTRAINT unique_name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Roles_20170121_013956) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `roles` CASCADE")
}
