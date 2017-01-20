package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UserMetas_20170121_014145 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UserMetas_20170121_014145{}
	m.Created = "20170121_014145"
	migration.Register("UserMetas_20170121_014145", m)
}

// Run the migrations
func (m *UserMetas_20170121_014145) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE `user_metas` ( "+
		"`id` Int( 32 ) AUTO_INCREMENT NOT NULL,"+
		"`user_id` Int( 32 ) NOT NULL,"+
		"`key` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,"+
		"`value` VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,"+
		"PRIMARY KEY ( `id` ),"+
		"CONSTRAINT `unique_id` UNIQUE( `id` ) ) "+
		"CHARACTER SET = utf8 "+
		"COLLATE = utf8_general_ci "+
		"ENGINE = InnoDB "+
		"AUTO_INCREMENT = 57;")
}

// Reverse the migrations
func (m *UserMetas_20170121_014145) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `user_metas` CASCADE")
}
