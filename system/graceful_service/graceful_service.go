package graceful_service

import (
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"syscall"
)

var (
	log = logging.MustGetLogger("graceful_service")
)

type GracefulService struct {
	cfg  *GracefulServiceConfig
	pool *GracefulServicePool
	done chan struct{}
}

func NewGracefulService(cfg *GracefulServiceConfig,
	hub *GracefulServicePool) (graceful *GracefulService) {
	graceful = &GracefulService{
		cfg:  cfg,
		pool: hub,
		done: make(chan struct{}, 1),
	}

	log.Info("Graceful shutdown service started")

	return
}

func (p GracefulService) Wait() {

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		<-gracefulStop
		p.pool.shutdown()
		p.done <- struct{}{}

	}()

	for {
		select {
		case <-p.done:
			log.Info("Shutdown")
			os.Exit(0)
		}
	}

	close(p.done)
	close(gracefulStop)
}

func (p GracefulService) Subscribe(client IGracefulClient) (id int) {
	id = p.pool.subscribe(client)
	return
}

func (p GracefulService) Unsubscribe(id int) {
	p.pool.unsubscribe(id)
	return
}


func (p GracefulService) Shutdown() {
	p.pool.shutdown()
	return
}
