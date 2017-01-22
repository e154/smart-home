package common

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"io/ioutil"
)

func GetScript(name string) (script []byte) {

	keys_path := beego.AppConfig.String("scripts_path")
	dir := filepath.Join("data", keys_path, name)

	if(beego.BConfig.RunMode == "dev") {
		dir = filepath.Join("../../", dir)
	}

	// Load sample key data
	var err error
	if script, err = ioutil.ReadFile(dir); err != nil {
		panic(err)
	}

	return
}

