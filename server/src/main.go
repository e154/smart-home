package main

import (
	"./api"
	"time"
)

func main() {
	api.Initialize()

	for ;; {
		time.Sleep(time.Second * 1)
	}
}
