package worker_manager

import (
	"log"
	"time"
)

// Singleton
var instantiated *WorkerManager = nil

func WorkerManagerPtr() *WorkerManager {
	return instantiated
}

type WorkerManager struct {

}

//TODO need normal cron like runner
func (wm *WorkerManager) Run(timer string, h func()) *WorkerManager {

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
	log.Println("Worker Manager initialize...")

	if instantiated == nil {
		instantiated = &WorkerManager{}
	}

	return
}