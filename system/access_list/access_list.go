package access_list

import (
	"io/ioutil"
	"encoding/json"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("access_list")
)

type AccessListService struct {
	List *AccessList
}

func NewAccessListService() *AccessListService {
	accessList := &AccessListService{}
	accessList.ReadConfig("./conf/access_list.json")
	return accessList
}

func (a *AccessListService) ReadConfig(path string) (err error) {

	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading config file")
		return
	}

	a.List = &AccessList{}
	err = json.Unmarshal(file, a.List)
	if err != nil {
		log.Fatal("Error: wrong format of config file")
		return
	}

	return
}
