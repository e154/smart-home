package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type MapImages_20170121_013529 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &MapImages_20170121_013529{}
	m.Created = "20170121_013529"
	migration.Register("MapImages_20170121_013529", m)
}

// Run the migrations
func (m *MapImages_20170121_013529) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE map_images (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	image_id Int( 32 ) NOT NULL,
	style Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_id UNIQUE( id ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *MapImages_20170121_013529) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `map_images` CASCADE")
}
