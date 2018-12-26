package use_case

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
	"errors"
)

func AddUser(params models.NewUser,
	adaptors *adaptors.Adaptors,
	currentUser *m.User) (ok bool, createdUser *m.User, errs []*validation.Error, err error) {

	// validation income request
	ok, errs = params.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	user := &m.User{}
	if err = copier.Copy(&user, &params); err != nil {
		return
	}

	if currentUser != nil {
		user.UserId.Scan(currentUser.Id)
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
	if len(errs) > 0 || !ok {
		return
	}

	var id int64
	if id, err = adaptors.User.Add(user); err != nil {
		return
	}

	createdUser, err = adaptors.User.GetById(id)

	return
}

func GetUserById(userId int64, adaptors *adaptors.Adaptors) (u *models.UserFullModel, err error) {

	var user *m.User
	if user, err = adaptors.User.GetById(userId); err != nil {
		return
	}

	// base model
	u = &models.UserFullModel{}
	copier.Copy(&u, &user)

	// parent model
	if user.User != nil {
		u.User = &models.UserByIdModelParent{}
		copier.Copy(&u.User, &user.User)
	}

	// meta
	u.Meta = make([]*models.UserByIdModelMeta, 0)
	for _, meta := range user.Meta {
		m := &models.UserByIdModelMeta{}
		copier.Copy(&m, &meta)
		u.Meta = append(u.Meta, m)
	}

	// history
	u.History = make([]*models.UserHistory, 0)
	for _, story := range user.History {
		s := &models.UserHistory{}
		copier.Copy(&s, &story)
		u.History = append(u.History, s)
	}

	// role
	u.Role = &models.Role{}
	copier.Copy(&u.Role, &user.Role)

	// image
	if user.Image != nil {
		u.Image = &models.Image{}
		copier.Copy(&u.Image, &user.Image)
	}

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

func GetUserList(limit, offset int, order, sortBy string, adaptors *adaptors.Adaptors) (items []*models.UserShotModel, total int64, err error) {

	var userList []*m.User
	if userList, total, err = adaptors.User.List(int64(limit), int64(offset), order, sortBy); err != nil {
		return
	}

	for _, user := range userList {
		item := &models.UserShotModel{}
		copier.Copy(&item , &user)

		// parent model
		if user.User != nil {
			item.User = &models.UserByIdModelParent{}
			copier.Copy(&item.User, &user.User)
		}

		// role
		item.Role = &models.Role{}
		copier.Copy(&item.Role, &user.Role)

		// image
		if user.Image != nil {
			item.Image = &models.Image{}
			copier.Copy(&item.Image, &user.Image)
		}

		items = append(items, item)
	}

	return
}

func UpdateUser(newParams *models.UpdateUser, adaptors *adaptors.Adaptors) (ok bool, errs []*validation.Error, err error) {

	ok, errs = newParams.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	var user *m.User
	if user, err = adaptors.User.GetById(newParams.Id); err != nil {
		return
	}

	if newParams.Password != newParams.PasswordRepeat {
		err = errors.New("bad passwords")
		return
	}

	copier.Copy(&user, &newParams)
	user.EncryptedPassword = common.Pwdhash(newParams.Password)

	err = adaptors.User.Update(user)

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