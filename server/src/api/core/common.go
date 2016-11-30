package core

import (
	"encoding/json"
	"log"
)

func debug(v interface{}) {
	j, _ := json.Marshal(v)
	log.Println(string(j))
}
