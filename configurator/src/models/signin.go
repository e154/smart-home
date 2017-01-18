package models

type Message struct {
	Status	string	`json:"status"`
	Message	string	`json:"message"`
}

type Signin struct {
	Token	string		`json:"token"`
	User	*User	`json:"current_user"`
}
