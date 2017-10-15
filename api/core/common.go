package core

import (
	"encoding/json"
	"log"
)

func debug(v interface{}) {
	j, _ := json.MarshalIndent(v, "", " ")
	log.Println(string(j))
}
