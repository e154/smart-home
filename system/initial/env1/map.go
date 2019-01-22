package env1

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

func addMaps(adaptors *adaptors.Adaptors,
	scripts map[string]*m.Script) (maps []*m.Map) {

	var err error

	// map 1
	m1 := &m.Map{
		Name:        "office1",
		Description: "офис на ул. Красный проспект, д.22",
		Options:     json.RawMessage(`{"zoom":1,"element_state_text":false,"element_option_text":false}`),
	}
	ok, _ := m1.Valid()
	So(ok, ShouldEqual, true)
	m1.Id, err = adaptors.Map.Add(m1)
	So(err, ShouldBeNil)

	return
}
