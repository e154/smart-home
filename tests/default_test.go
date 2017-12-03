package test

import (
	"os"
	"fmt"
	"testing"
	"net/http"
	"encoding/json"
	"path/filepath"
	"net/http/httptest"
	"github.com/astaxie/beego"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/DATA-DOG/godog"
	"github.com/e154/smart-home/database"
	server "github.com/e154/smart-home/api"
	"github.com/e154/smart-home/api/core"
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
	resp 	*httptest.ResponseRecorder
}

func (a *apiFeature) resetResponse(interface{}) {
	a.resp = httptest.NewRecorder()
}

func (a *apiFeature) iSendrequestTo(method, endpoint string) (err error) {

	uri := fmt.Sprintf("%s%s",a.basePath(), endpoint)
	//fmt.Println("url:", u)

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return
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

func (a *apiFeature) DropDb() error {

	// stop server
	core.CorePtr().Stop()

	// drop test database
	config, _ := beego.AppConfig.GetSection("test")
	database.DropDb(config["db_name"])

	database.Migration(database.GetDbConfig(true))

	// run server
	core.CorePtr().Run()

	return nil
}

func FeatureContext(s *godog.Suite) {
	api = &apiFeature{}

	s.BeforeScenario(func(interface{}) {
		api.resetResponse(nil)
	})

	s.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendrequestTo)
	s.Step(`^database is clean`, api.DropDb)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
}