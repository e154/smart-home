package models

type AccessItem struct {
	Actions		[]string	`json:"actions"`
	Method		string		`json:"method"`
	Description	string		`json:"description"`
	RoleName	string		`json:"role_name"`
}

type AccessLevels map[string]AccessItem
func NewAccessLevels() AccessLevels {
	return  make(map[string]AccessItem)
}

type AccessList map[string]AccessLevels
func NewAccessList() AccessList {
	return make(map[string]AccessLevels)
}
