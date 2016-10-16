package settings

import (
	"os"
	"time"
	"github.com/astaxie/beego/config"
	"fmt"
	"strconv"
	"io/ioutil"
	"path/filepath"
	"log"
)

const (
	CONF_NAME string = "node.conf"
	APP_NAME string = "node"
	APP_MAJOR = 0
	APP_MINOR = 1
	APP_PATCH = 0
	permMode os.FileMode = 0666
)

// Singleton
var instantiated *Settings = nil

func SettingsPtr() *Settings {
	if instantiated == nil {
		instantiated = new(Settings)
	}
	return instantiated
}

type Settings struct {
	IP		string
	Port		int
	DeviceList	[]string
	StartTime	time.Time
	UpTime		time.Duration
	cfg 		config.IniConfigContainer
	dir		string
}

func (s *Settings) Init() *Settings {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	s.StartTime = time.Now()
	s.dir = fmt.Sprintf("%s/", dir)
	s.IP = "127.0.0.1"
	s.Port = 8888

	s.Load()

	return s
}

func (s *Settings) Load() (*Settings, error) {

	if _, err := os.Stat(s.dir + CONF_NAME); os.IsNotExist(err) {
		return s.Save()
	}

	// read config file
	cfg, err := config.NewConfig("ini", s.dir + CONF_NAME)
	if err != nil {
		return s, err
	}

	if cfg.String("app_version") != s.AppVresion() {
		s.Save()
		return s.Load()
	}

	s.IP = cfg.String("ip")
	s.Port, _ = cfg.Int("port")

	return s, err
}

func (s *Settings) Save() (*Settings, error) {

	if _, err := os.Stat(s.dir + CONF_NAME); os.IsNotExist(err) {
		ioutil.WriteFile(s.dir + CONF_NAME, []byte{}, permMode)
	}

	cfg, err := config.NewConfig("ini", s.dir + CONF_NAME)
	if err != nil {
		return s, err
	}

	cfg.Set("app_version", s.AppVresion())
	cfg.Set("ip", s.IP)
	cfg.Set("port", strconv.Itoa(s.Port))

	if err := cfg.SaveConfigFile(s.dir + CONF_NAME); err != nil {
		fmt.Printf("err with create conf file: %s\n", s.dir + CONF_NAME)
		return s, err
	}

	return s, nil
}

func (s *Settings) AppVresion() string {
	return fmt.Sprintf("%d.%d.%d", APP_MAJOR, APP_MINOR, APP_PATCH)
}

func init() {
	instantiated = new(Settings)
}