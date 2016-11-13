package crontab

import (
	"log"
	"time"
)

// Singleton
var instantiated *Crontab = nil

func WorkerManagerPtr() *Crontab {
	return instantiated
}

type Crontab struct {

}

//TODO need normal cron like runner
func (wm *Crontab) Run(timer string, h func()) *Crontab {

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
			// do stuff
			h()

			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	return wm
}

func Initialize() (err error) {
	log.Println("Crontab initialize...")

	if instantiated == nil {
		instantiated = &Crontab{}
	}

	return
}