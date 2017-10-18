package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type AddUserLanguage_20171018_191650 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AddUserLanguage_20171018_191650{}
	m.Created = "20171018_191650"
	migration.Register("AddUserLanguage_20171018_191650", m)
}

// Run the migrations
func (m *AddUserLanguage_20171018_191650) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("ALTER TABLE `users` ADD COLUMN `lang` VarChar( 255 ) NOT NULL DEFAULT 'en';")
}

// Reverse the migrations
func (m *AddUserLanguage_20171018_191650) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("ALTER TABLE `users` DROP COLUMN `lang`;")
}
