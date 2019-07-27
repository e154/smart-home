package models

import (
	"encoding/json"
	"time"
)

type Variable struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Autoload  bool      `json:"autoload"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewVariable(name string) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) GetObj(obj interface{}) (err error) {
	err = json.Unmarshal([]byte(v.Value), obj)
	return
}

func (v *Variable) SetObj(obj interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(obj); err != nil {
		return
	}
	v.Value = string(b)
	return
}
