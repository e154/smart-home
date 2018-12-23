package use_case

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
)

func AddUser(params models.NewUser, adaptors *adaptors.Adaptors) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation income request
	ok, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	user := &m.User{}
	if err = copier.Copy(&user, &params); err != nil {
		return
	}

	if params.Password == params.PasswordRepeat {
		user.EncryptedPassword = common.Pwdhash(params.Password)
	}

	if params.Avatar != nil && params.Avatar.Id != 0 {
		user.ImageId.Scan(params.Avatar.Id)
	}

	if params.Role != nil {
		user.RoleName = params.Role.Name
	}

	if params.Meta != nil && len(params.Meta) > 0 {
		for _, rMeta := range params.Meta {
			meta := &m.UserMeta{}
			if err = copier.Copy(&meta, &rMeta); err != nil {
				return
			}
			user.Meta = append(user.Meta, meta)
		}
	}

	// check user status
	switch user.Status {
	case "active", "blocked":
	default:
		user.Status = "blocked"
	}

	// validation user model
	ok, errs = user.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = adaptors.User.Add(user); err != nil {
		return
	}

	user.Id = id

	return
}
