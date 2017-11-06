package variable

import (
	"sync"
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/models"
)

// Singleton
var instantiated *variable = nil

func Load(name string) (string, error) {
	return instantiated.Read(name)
}

func Get(name string) (value string, ok bool) {
	instantiated.Lock()
	defer instantiated.Unlock()

	value, ok = instantiated.pull[name]

	return
}

func Set(name, value string) error {
	instantiated.Lock()
	defer instantiated.Unlock()

	instantiated.pull[name] = value

	return instantiated.Save(name, value)
}

func Delete(name string) error {
	instantiated.Lock()
	defer instantiated.Unlock()

	return instantiated.Delete(name)
}

type variable struct {
	sync.Mutex
	pull	map[string]string
}

func (v *variable) Init() {
	v.pull = make(map[string]string)

	variables, err := models.GetAllVariable()
	if err != nil {
		log.Error(err.Error())
	}

	defer v.Unlock()
	v.Lock()
	for _, value := range variables {
		v.pull[value.Name] = value.Value
	}
}

func (v *variable) Read(name string) (string, error) {

	var variable *models.Variable
	var err error

	if variable, err = models.GetVariableByName(name); err != nil {
		return "", err
	}

	return variable.Value, nil
}

func (v *variable) Save(name, value string) error {

	variable := &models.Variable{
		Name: name,
		Value: value,
		Autoload: "yes",
	}

	return  models.InsertOrUpdateVariableByName(variable)
}

func (v *variable) Delete(name string) error {

	if _, ok := v.pull[name]; ok {
		delete(v.pull, name)
	}

	return models.DeleteVariable(name)
}

func Initialize() {
	log.Info("Settings initialize...")

	if instantiated == nil {
		instantiated = new(variable)
		instantiated.Init()
	}

	return
}