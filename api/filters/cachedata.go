package filters

import (
	"github.com/e154/smart-home/api/models"
)

type cacheData struct {
	user	*models.User
	access_list	models.AccessList
}