package mqtt

import (
	"fmt"
	"github.com/surgemq/surgemq/auth"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/uuid"
)

var ErrBadLoginOrPassword = fmt.Errorf("bad login or password")

type Authenticator struct {
	adaptors *adaptors.Adaptors
	name     string
	uuid     uuid.UUID
}

func NewAuthenticator(adaptors *adaptors.Adaptors) *Authenticator {
	a := &Authenticator{
		name:     "base",
		adaptors: adaptors,
		uuid:     uuid.NewV4(),
	}
	a.Register()
	return a
}

func (a *Authenticator) Authenticate(login string, pass interface{}) (err error) {

	//fmt.Println(login, pass)

	password, ok := pass.(string)
	if !ok || password == "" {
		err = ErrBadLoginOrPassword
	}

	if login == "local" && pass == a.uuid.String() {
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

func (a *Authenticator) Register() {
	auth.Register(a.name, a)
	return
}

func (a Authenticator) Name() string {
	return a.name
}

func (a Authenticator) LocalClientUuid() string {
	return a.uuid.String()
}
