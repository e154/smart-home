package test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

type signInList struct {
	Username 	string
	Password 	string
	RespCode 	int
}

func TestSignin(t *testing.T)  {
	api.resetResponse(nil)

	var states = []struct {
		Username 	string
		Password 	string
		RespCode 	int
	}{
		{"test@e154.ru1", "testtest2", 401},
		{"test@e154.ru1", "testtest", 401},
		{"test@e154.ru", "testtest2", 403},
		{"test@e154.ru", "testtest", 201},
	}

	for _, v := range states {
		api.Authorization(v.Username, v.Password)

		Convey("Subject: Test Station Endpoint\n", t, func() {
			Convey("Status Code", func() {
				So(api.resp.Code, ShouldEqual, v.RespCode)
			})
			Convey("The Result Should Not Be Empty", func() {
				So(api.resp.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})

		api.resetResponse(nil)
	}
}

func TestSignout(t *testing.T)  {

	api.resetResponse(nil)

	var states = []struct {
		Username 	string
		Password 	string
		RespCode 	int
	}{
		{"test@e154.ru", "testtest", 201},
		{"test@e154.ru", "testtest2", 403},
	}

	for _, v := range states {
		api.Authorization(v.Username, v.Password)

		api.iSendrequestTo("POST", "/api/v1/signout")

		Convey("Subject: Test Station Endpoint\n", t, func() {
			Convey("Status Code", func() {
				So(api.resp.Code, ShouldEqual, v.RespCode)
			})
			Convey("The Result Should Not Be Empty", func() {
				So(api.resp.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})

		api.resetResponse(nil)
	}
}

func TestRecovery(t *testing.T)  {}

func TestReset(t *testing.T)  {}

func TestAccessList(t *testing.T)  {}
