package scripts

import (
	"fmt"
	"errors"
	"github.com/op/go-logging"
	"github.com/e154/smart-home/system/config"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/scripts/bind"
)

var (
	log = logging.MustGetLogger("scripts")
)

type ScriptService struct {
	cfg  *config.AppConfig
	pull *Pull
}

func NewScriptService(cfg *config.AppConfig) (service *ScriptService) {

	pull := &Pull{
		functions:  make(map[string]interface{}),
		structures: make(map[string]interface{}),
	}

	service = &ScriptService{cfg: cfg, pull: pull}
	//service.PushStruct("Notifr", &bind.NotifrBind{})
	service.PushStruct("Log", &bind.LogBind{})
	service.PushFunctions("ExecuteSync", bind.ExecuteSync)
	service.PushFunctions("ExecuteAsync", bind.ExecuteAsync)
	return service
}

func (service ScriptService) NewEngine(s *m.Script) (engine *Engine, err error) {

	engine = &Engine{
		model: s,
		buf:   make([]string, 0),
		pull:  service.pull,
	}

	switch s.Lang {
	case ScriptLangTs, ScriptLangCoffee, ScriptLangJavascript:
		engine.script = &Javascript{engine: engine}
	default:
		err = errors.New(fmt.Sprintf("undefined language %s", s.Lang))
		return
	}

	//if err == nil {
	//	log.Infof("Add script: %s (%s)", s.Name, s.Lang)
	//}

	engine.script.Init()

	return
}

func (service *ScriptService) PushStruct(name string, s interface{}) {
	service.pull.Lock()
	defer service.pull.Unlock()

	service.pull.structures[name] = s
}

func (service *ScriptService) PushFunctions(name string, s interface{}) {
	service.pull.Lock()
	defer service.pull.Unlock()

	//fmt.Println("push function")

	service.pull.functions[name] = s
}
