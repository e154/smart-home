package models

import "time"

type Role struct {
	Name		string			`json:"name"`
	Description	string			`json:"description"`
	Parent		*Role			`json:"parent"`
	Children	[]*Role			`json:"children"`
	AccessList	map[string][]string	`json:"access_list"`
	Created_at	time.Time		`json:"created_at"`
	Update_at	time.Time		`json:"update_at"`
}
