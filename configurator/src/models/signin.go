package models

type Signin struct {
	Token	string		`json:"token"`
	User	*User	`json:"current_user"`
}
