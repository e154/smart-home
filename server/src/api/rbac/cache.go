package rbac

import (
	"../models"
)

type cacheData struct {
	user	*models.User
	access_list	models.AccessList
}