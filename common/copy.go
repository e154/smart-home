package common

import (
	"github.com/jinzhu/copier"
	"encoding/json"
	"bytes"
	"encoding/gob"
)

type CopyEngine string

const (
	JsonEngine = CopyEngine("json")
	GobEngine  = CopyEngine("gob")
)

func gobCopy(to, from interface{}) (err error) {
	buff := new(bytes.Buffer)
	if err = gob.NewEncoder(buff).Encode(from); err != nil {
		return
	}
	err = gob.NewDecoder(buff).Decode(to)
	return
}

func jsonCopy(to, from interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(from); err != nil {
		return
	}
	err = json.Unmarshal(b, to)
	return
}

func Copy(to, from interface{}, params ...CopyEngine) (err error) {

	if len(params) == 0 {
		err = copier.Copy(to, from)
		return
	}

	switch params[0] {
	case JsonEngine:
		err = jsonCopy(to, from)
	case GobEngine:
		err = gobCopy(to, from)
	}

	return
}
