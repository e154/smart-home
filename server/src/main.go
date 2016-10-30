package main

import (
	"github.com/astaxie/beego"
	"./api"
	"time"
	"net"
	"bytes"
	"io"
	"encoding/json"
	"fmt"
)

func main() {
	api.Initialize()
	beego.Info("Starting....")
	go beego.Run()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
