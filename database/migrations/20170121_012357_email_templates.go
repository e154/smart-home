package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type EmailTemplates_20170121_012357 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &EmailTemplates_20170121_012357{}
	m.Created = "20170121_012357"
	migration.Register("EmailTemplates_20170121_012357", m)
}

// Run the migrations
func (m *EmailTemplates_20170121_012357) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE email_templates (
	name VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'Название',
	description VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'Описание',
	content LongText CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'Содержимое',
	status VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'active' COMMENT 'active, unactive',
	type VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'item' COMMENT 'item, template',
	parent VarChar( 64 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	created_at DateTime NOT NULL,
	updated_at DateTime NOT NULL,
	PRIMARY KEY ( name ),
	CONSTRAINT name UNIQUE( name ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *EmailTemplates_20170121_012357) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `email_templates` CASCADE")
}
