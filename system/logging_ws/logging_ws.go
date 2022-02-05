package logging_ws

import (
	"context"
	"encoding/json"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/stream"
	"go.uber.org/fx"
)

var (
	log = common.MustGetLogger("logging_ws")
)

type LoggingWs struct {
	logging *logging.Logging
	stream  *stream.Stream
}

func NewLogWsSaver(lc fx.Lifecycle,
	logging *logging.Logging,
	stream *stream.Stream) *LoggingWs {

	l := &LoggingWs{
		logging: logging,
		stream:  stream,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return l.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return l.Shutdown(ctx)
		},
	})

	return l
}

func (l *LoggingWs) Start(ctx context.Context) (err error) {
	log.Info("start")
	l.logging.SetWsSaver(l)
	return
}

func (l *LoggingWs) Shutdown(ctx context.Context) (err error) {
	log.Info("shutdown")
	l.logging.SetWsSaver(nil)
	return
}

// Save ...
func (l *LoggingWs) Save(log m.Log) {
	b, _ := json.Marshal(log)
	l.stream.Broadcast("log", b)
}
