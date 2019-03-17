package models

// swagger:model
type AuthSignInResponse struct {
	CurrentUser *CurrentUser `json:"current_user"`
	AccessToken string       `json:"access_token"`
}
