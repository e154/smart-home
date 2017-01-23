package main

import (
	"github.com/e154/smart-home/api"
	"time"
)

func main() {
	api.Initialize()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
