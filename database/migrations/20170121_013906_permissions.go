package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Permissions_20170121_013906 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Permissions_20170121_013906{}
	m.Created = "20170121_013906"
	migration.Register("Permissions_20170121_013906", m)
}

// Run the migrations
func (m *Permissions_20170121_013906) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE permissions (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	role_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	package_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	level_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Permissions_20170121_013906) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `permissions` CASCADE")
}
