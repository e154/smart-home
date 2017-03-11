package common

import (
	"math/rand"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	crypto_rand "crypto/rand"

	"github.com/astaxie/beego/validation"
	"github.com/fiam/gounidecode/unidecode"
	"strings"
	"os"
	"os/exec"
	"github.com/astaxie/beego"
	"time"
	"path/filepath"
	"unicode/utf8"
)

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func ErrorParse(valid validation.Validation) string {
	var msg string
	for _, err := range valid.Errors {
		msg += fmt.Sprintf( "%s: %s\r", err.Key, err.Message)
	}

	return msg
}

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//password hash function
func Pwdhash(str string) string {
	return Strtomd5(str)
}

func StringsToJson(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

func RandStr(strSize int, randType string) string {

	var dictionary string

	switch randType {
	case "alphanum":
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "alpha":
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "number":
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	crypto_rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func GenUrl(s string) (url string) {
	return strings.Replace(strings.ToLower(unidecode.Unidecode(s)), " ", "_", -1)
}

func GetFileSize(name string) (int64, error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func GenUuid() (string, error) {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		beego.Warn(err.Error())
		return "", err
	}
	return string(uuid), nil
}

func GetFullPath(name string) string {
	data_dir := beego.AppConfig.String("data_dir")
	file_storage_path := beego.AppConfig.String("file_storage_path")
	dir := filepath.Join(data_dir, file_storage_path)

	depth, err := beego.AppConfig.Int("file_storage_depth")
	if err != nil {
		depth = 3
	}

	for i := 0; i < depth; i++ {
		dir = filepath.Join(dir, name[i*3:(i+1)*3])
	}

	return dir
}

func GetLinkPath(name string) string {

	dir := "/static"

	count := utf8.RuneCountInString(name)
	if count < 9 {
		return filepath.Join(dir, name)
	}

	depth, err := beego.AppConfig.Int("file_storage_depth")
	if err != nil {
		depth = 3
	}

	for i := 0; i < depth; i++ {
		dir = filepath.Join(dir, name[i*3:(i+1)*3])
	}

	return filepath.Join(dir, name)
}