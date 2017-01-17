package models

import (
	"time"
	"encoding/gob"
)

type User struct {
	Id                     	int64                  	`json:"id"`
	Nickname               	string                	`json:"nickname"`
	FirstName             	string                	`json:"first_name"`
	LastName              	string                	`json:"last_name"`

	Email                  	string                	`json:"email"`
	History                	string                	`json:"history"`
	Status                 	string                	`json:"status"`

	Avatar                 	*Image                	`json:"avatar"`
	SignInCount          	int64                   `json:"sign_in_count"`
	CurrentSignInIp     	string                	`json:"current_sign_in_ip"`
	LastSignInIp        	string                	`json:"last_sign_in_ip"`

	CreatedBy	       	*User                   `json:"created_by"`
	Role                   	*Role                	`json:"role"`
	Meta			[]*UserMeta		`json:"meta"`

	CurrentSignInAt     	time.Time        	`json:"current_sign_in_at"`
	LastSignInAt        	time.Time        	`json:"last_sign_in_at"`
	Created_at            	time.Time        	`json:"created_at"`
	Update_at            	time.Time       	`json:"update_at"`


}

func init() {
	gob.Register(new(User))
}
