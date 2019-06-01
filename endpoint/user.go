package endpoint

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"errors"
)

type UserEndpoint struct {
	*CommonEndpoint
}

func NewUserEndpoint(common *CommonEndpoint) *UserEndpoint {
	return &UserEndpoint{
		CommonEndpoint: common,
	}
}

func (n *UserEndpoint) Add(params *m.User,
	currentUser *m.User) (result *m.User, errs []*validation.Error, err error) {

	user := &m.User{}
	if err = common.Copy(&user, &params); err != nil {
		return
	}

	_, errs = user.Valid()
	if len(errs) > 0 {
		return
	}

	if currentUser != nil {
		user.UserId.Scan(currentUser.Id)
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
	if id, err = n.adaptors.User.Add(user); err != nil {
		return
	}

	result, err = n.GetById(id)

	return
}

func (n *UserEndpoint) GetById(userId int64) (result *m.User, err error) {

	result, err = n.adaptors.User.GetById(userId)

	return
}

func (n *UserEndpoint) Delete(userId int64) (err error) {

	var user *m.User
	if user, err = n.adaptors.User.GetById(userId); err != nil {
		return
	}

	if user.Role.Name == "admin" {
		err = errors.New("admin is main user")
		return
	}

	err = n.adaptors.User.Delete(user.Id)

	return
}

func (n *UserEndpoint) GetList(limit, offset int, order, sortBy string) (result []*m.User, total int64, err error) {

	result, total, err = n.adaptors.User.List(int64(limit), int64(offset), order, sortBy)

	return
}

func (n *UserEndpoint) Update(params *m.User) (result *m.User, errs []*validation.Error, err error) {

	var user *m.User
	if user, err = n.adaptors.User.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&user, &params, common.JsonEngine)

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

	if err = n.adaptors.User.Update(user); err != nil {
		return
	}

	result, err = n.GetById(user.Id)

	return
}

func (n *UserEndpoint) UpdateStatus(userId int64, newStatus string) (err error) {

	var user *m.User
	if user, err = n.adaptors.User.GetById(userId); err != nil {
		return
	}

	user.Status = newStatus

	// check user status
	switch user.Status {
	case "active", "blocked":
	default:
		user.Status = "blocked"
	}

	err = n.adaptors.User.Update(user)

	return
}
