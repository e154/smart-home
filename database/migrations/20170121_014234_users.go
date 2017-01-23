package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Users_20170121_014234 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Users_20170121_014234{}
	m.Created = "20170121_014234"
	migration.Register("Users_20170121_014234", m)
}

// Run the migrations
func (m *Users_20170121_014234) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE users (
	id Int( 32 ) AUTO_INCREMENT NOT NULL,
	nickname VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	first_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	last_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	encrypted_password VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	email VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	history Text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	status VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	reset_password_token VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	authentication_token VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	image_id Int( 255 ) NULL,
	sign_in_count Int( 32 ) NOT NULL,
	current_sign_in_ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	last_sign_in_ip VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NULL,
	user_id Int( 32 ) NULL,
	role_name VarChar( 255 ) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
	reset_password_sent_at DateTime NULL,
	current_sign_in_at DateTime NULL,
	last_sign_in_at DateTime NULL,
	created_at DateTime NOT NULL,
	update_at DateTime NOT NULL,
	deleted DateTime NULL,
	PRIMARY KEY ( id ),
	CONSTRAINT unique_email UNIQUE( email ),
	CONSTRAINT unique_id UNIQUE( id ),
	CONSTRAINT unique_nickname UNIQUE( nickname ) )
	CHARACTER SET = utf8
	COLLATE = utf8_general_ci
	ENGINE = InnoDB;`)
}

// Reverse the migrations
func (m *Users_20170121_014234) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS `users` CASCADE")
}
