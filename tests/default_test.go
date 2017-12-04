package test

import (
	"os"
	"fmt"
	"bytes"
	"testing"
	"net/http"
	"encoding/json"
	"path/filepath"
	"encoding/base64"
	"net/http/httptest"
	"github.com/astaxie/beego"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/DATA-DOG/godog"
	"github.com/e154/smart-home/database"
	server "github.com/e154/smart-home/api"
	"github.com/e154/smart-home/api/core"
	"github.com/astaxie/beego/orm"
)

var (
	httpaddr	string
	httpport	string
	api 		*apiFeature
)

func init() {
	apppath := filepath.Join(os.Getenv("PWD"), "..")
	beego.Info("init:", apppath)

	beego.TestBeegoInit(apppath)

	httpport = beego.AppConfig.String("httpport")
	httpaddr = beego.AppConfig.String("httpaddr")

	beego.AppConfig.Set("orm_debug", "false")

	database.Initialize(true)

	// drop test database
	config, _ := beego.AppConfig.GetSection("test")
	database.DropDb(config["db_name"])

	// run migration
	database.Migration(database.GetDbConfig(true))

	server.Initialize(true)
}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "progress",
		Paths:     []string{"tests/features"},
		//Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}


type apiFeature struct {
	resp 			*httptest.ResponseRecorder
	headParams 		map[string]string
	jsonStr 		[]byte
	access_token	string
}

func (a *apiFeature) resetResponse(interface{}) {
	a.resp = httptest.NewRecorder()
}

func (a *apiFeature) iSendrequestTo(method, endpoint string) (err error) {

	uri := fmt.Sprintf("%s%s",a.basePath(), endpoint)
	//fmt.Println("url:", uri)

	req, _ := http.NewRequest(method, uri, bytes.NewBuffer(a.jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = endpoint

	if a.headParams != nil {
		for k, v := range a.headParams {
			req.Header.Set(k, v)
		}
	}

	if a.access_token != "" {
		req.Header.Set("access_token", a.access_token)
	}

	// handle panic
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	beego.BeeApp.Handlers.ServeHTTP(a.resp, req)

	return
}

func (a *apiFeature) basePath() string {

	return fmt.Sprintf("http://%s:%s", httpaddr, httpport)
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}
	return nil
}

func (a *apiFeature) theResponseShouldMatchJSON(body *gherkin.DocString) (err error) {
	var expected, actual []byte
	var exp, act interface{}

	// re-encode expected response
	if err = json.Unmarshal([]byte(body.Content), &exp); err != nil {
		return
	}

	if expected, err = json.MarshalIndent(exp, "", "  "); err != nil {
		return
	}

	// re-encode actual response too
	if err = json.Unmarshal(a.resp.Body.Bytes(), &act); err != nil {
		return
	}
	if actual, err = json.MarshalIndent(act, "", "  "); err != nil {
		return
	}

	// the matching may be adapted per different requirements.
	if len(actual) != len(expected) {
		return fmt.Errorf(
			"expected json length: %d does not match actual: %d:\n%s",
			len(expected),
			len(actual),
			string(actual),
		)
	}

	for i, b := range actual {
		if b != expected[i] {
			return fmt.Errorf(
				"expected JSON does not match actual, showing up to last matched character:\n%s",
				string(actual[:i+1]),
			)
		}
	}
	return
}

func (a *apiFeature) Authorization(name, password string) (err error) {

	var auth string = base64.StdEncoding.EncodeToString([]byte(name + ":" + password))

	a.headParams["Authorization"] = "Basic " + auth

	a.jsonStr = []byte{}

	if err = a.iSendrequestTo("GET", "/api/v1/signin"); err != nil {
		return
	}

	act := make(map[string]interface{})
	if err = json.Unmarshal(a.resp.Body.Bytes(), &act); err != nil {
		return
	}

	if _, ok := act["access_token"]; ok {
		a.access_token = act["access_token"].(string)
	}

	return
}

func (a *apiFeature) iFinishingTheSession() (err error) {

	if a.access_token == "" {
		return
	}

	err = a.iSendrequestTo("POST", "/api/v1/signout")

	return
}

func (a *apiFeature) DropDb() error {

	// stop server
	core.CorePtr().Stop()

	// drop test database
	config, _ := beego.AppConfig.GetSection("test")
	database.DropDb(config["db_name"])

	database.Migration(database.GetDbConfig(true))

	// run server
	core.CorePtr().Run()

	a.DbPrepare()

	return nil
}

func (a *apiFeature) DbPrepare() {

	o := orm.NewOrm()

	o.Raw(`INSERT INTO users ( created_at, update_at, email, encrypted_password, first_name,
	id, lang, last_name, nickname, role_name, status, history, reset_password_token, authentication_token, sign_in_count, current_sign_in_ip)
	VALUES ( '2017-12-04 15:08:56', '2017-12-04 15:08:56', 'test@e154.ru', '05a671c66aefea124cc08b76ea6d30bb',
	 'John', 3, 'en', 'Doe', 'test', 'demo', 'active', '[]', '', '', 0, '');`).Exec()

}

func FeatureContext(s *godog.Suite) {
	api = &apiFeature{
		headParams: make(map[string]string),
	}

	api.DbPrepare()

	s.BeforeScenario(func(interface{}) {
		api.resetResponse(nil)
	})

	s.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendrequestTo)
	s.Step(`^database is clean`, api.DropDb)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
	s.Step(`^I authorization with user "([^"]*)" and password "([^"]*)"$`, api.Authorization)
	s.Step(`^I finishing the session$`, api.iFinishingTheSession)
}