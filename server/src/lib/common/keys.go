package common

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/astaxie/beego"
	"path/filepath"
	"io/ioutil"
)

func GetKey(name string) (keyData []byte) {

	keys_path := beego.AppConfig.String("keys_path")
	dir := filepath.Join("data", keys_path, name)

	if(beego.BConfig.RunMode == "dev") {
		dir = filepath.Join("../../", dir)
	}

	// Load sample key data
	var err error
	if keyData, err = ioutil.ReadFile(dir); err != nil {
		panic(err)
	}

	return
}

func GetHmacToken(data map[string]interface{}, key []byte) (tokenString string, err error){

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))

	// Sign and get the complete encoded token as a string using the secret
	if tokenString, err = token.SignedString(key); err != nil {
		return
	}

	return
}

func ParseHmacToken(tokenString string, key []byte) (jwt.MapClaims, error) {

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return key, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}