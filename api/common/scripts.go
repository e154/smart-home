package common

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"io/ioutil"
)

func GetScript(name string) (script []byte) {

	keys_path := beego.AppConfig.String("scripts_path")
	data_dir := beego.AppConfig.String("data_dir")
	dir := filepath.Join(data_dir, keys_path, name)

	// Load sample key data
	var err error
	if script, err = ioutil.ReadFile(dir); err != nil {
		panic(err)
	}

	return
}

