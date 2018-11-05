package config

import (
	"io/ioutil"
	"log"
	"encoding/json"
)

const path = "conf/config.json"

func ReadConfig() (conf *AppConfig, err error) {
	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading config file")
		return
	}
	conf = &AppConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Error: wrong format of config file")
		return
	}

	return
}
