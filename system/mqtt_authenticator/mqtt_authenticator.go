package mqtt_authenticator

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("mqtt_authenticator")
)

var ErrBadLoginOrPassword = fmt.Errorf("bad login or password")

type Authenticator struct {
	adaptors *adaptors.Adaptors
	name     string
	login    string
	password string
}

func NewAuthenticator(adaptors *adaptors.Adaptors) *Authenticator {
	a := &Authenticator{
		adaptors: adaptors,
		name:     "base",
		login:    "local",
		password: uuid.NewV4().String(),
	}
	return a
}

func (a *Authenticator) Authenticate(login string, pass interface{}) (err error) {

	log.Infof("login: %v, pass: %v", login, pass)

	password, ok := pass.(string)
	if !ok || password == "" {
		err = ErrBadLoginOrPassword
	}

	if login == a.login && pass == a.password {
		return
	}

	var node *m.Node
	if node, err = a.adaptors.Node.GetByLogin(login); err != nil {
		return
	}

	if ok := common.CheckPasswordHash(password, node.EncryptedPassword); !ok {
		err = ErrBadLoginOrPassword
	}

	return
}

func (a Authenticator) Name() string {
	return a.name
}

func (a Authenticator) Password() string {
	return a.password
}

func (a Authenticator) Login() string {
	return a.login
}
