package use_case

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
	"errors"
)

func AddUser(params models.NewUser,
	adaptors *adaptors.Adaptors,
	currentUser *m.User) (result *models.UserFull, errs []*validation.Error, err error) {

	user := &m.User{}
	if err = common.Copy(&user, &params); err != nil {
		return
	}

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if currentUser != nil {
		user.UserId.Scan(currentUser.Id)
	}

	if params.Password == params.PasswordRepeat {
		user.EncryptedPassword = common.Pwdhash(params.Password)
	}

	if params.Image != nil && params.Image.Id != 0 {
		user.ImageId.Scan(params.Image.Id)
	}

	if params.Role != nil {
		user.RoleName = params.Role.Name
	}

	if params.Meta != nil && len(params.Meta) > 0 {
		for _, rMeta := range params.Meta {
			meta := &m.UserMeta{}
			if err = common.Copy(&meta, &rMeta); err != nil {
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
	_, errs = user.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = adaptors.User.Add(user); err != nil {
		return
	}

	result, err = GetUserById(id, adaptors)

	return
}

func GetUserById(userId int64, adaptors *adaptors.Adaptors) (result *models.UserFull, err error) {

	var user *m.User
	if user, err = adaptors.User.GetById(userId); err != nil {
		return
	}

	// base model
	result = &models.UserFull{}
	err = common.Copy(&result, &user, common.JsonEngine)

	return
}

func DeleteUserById(userId int64, adaptors *adaptors.Adaptors) (err error) {

	var user *m.User
	if user, err = adaptors.User.GetById(userId); err != nil {
		return
	}

	if user.Role.Name == "admin" {
		err = errors.New("admin is main user")
		return
	}

	err = adaptors.User.Delete(user.Id)

	return
}

func GetUserList(limit, offset int, order, sortBy string, adaptors *adaptors.Adaptors) (items []*models.UserShot, total int64, err error) {

	var userList []*m.User
	if userList, total, err = adaptors.User.List(int64(limit), int64(offset), order, sortBy); err != nil {
		return
	}

	for _, user := range userList {
		item := &models.UserShot{}
		common.Copy(&item, &user)

		// parent model
		if user.User != nil {
			item.User = &models.UserByIdModelParent{}
			common.Copy(&item.User, &user.User)
		}

		// role
		item.Role = &models.Role{}
		common.Copy(&item.Role, &user.Role)

		// image
		if user.Image != nil {
			item.Image = &models.Image{}
			common.Copy(&item.Image, &user.Image)
		}

		items = append(items, item)
	}

	return
}

func UpdateUser(params *models.UpdateUser, adaptors *adaptors.Adaptors) (result *models.UserFull, errs []*validation.Error, err error) {

	var user *m.User
	if user, err = adaptors.User.GetById(params.Id); err != nil {
		return
	}

	if params.Password != "" && params.Password != params.PasswordRepeat {
		err = errors.New("bad passwords")
		return
	}

	common.Copy(&user, &params, common.JsonEngine)
	user.EncryptedPassword = common.Pwdhash(params.Password)

	if params.Image != nil && params.Image.Id != 0 {
		user.ImageId.Scan(params.Image.Id)
	}

	if params.Role != nil {
		user.RoleName = params.Role.Name
	}

	if params.Meta != nil && len(params.Meta) > 0 {
		for _, rMeta := range params.Meta {
			meta := &m.UserMeta{}
			if err = common.Copy(&meta, &rMeta); err != nil {
				return
			}
			user.Meta = append(user.Meta, meta)
		}
	}

	_, errs = user.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.User.Update(user); err != nil {
		return
	}

	result, err = GetUserById(user.Id, adaptors)

	return
}

func UpdateStatus(userId int64, newStatus string, adaptors *adaptors.Adaptors) (err error) {

	var user *m.User
	if user, err = adaptors.User.GetById(userId); err != nil {
		return
	}

	user.Status = newStatus

	// check user status
	switch user.Status {
	case "active", "blocked":
	default:
		user.Status = "blocked"
	}

	err = adaptors.User.Update(user)

	return
}
