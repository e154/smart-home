package telemetry

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/e154/smart-home/common"
)

type StatusCode string

const (
	Ok    = StatusCode("Ok")
	Error = StatusCode("Error")
)

type Span struct {
	traceName   string
	name        string
	Ctx         context.Context
	started     time.Time
	ended       *time.Time
	Level       int
	Num         int
	status      StatusCode
	description string
	attributes  map[string]string
}

func Start(ctx context.Context, name string) (context.Context, *Span) {
	var level = 1
	var num = 1
	if p, ok := SpanFromContext(ctx); ok {
		num += p.Num
		level = p.Level
		if p.ended == nil {
			level += 1
		}
	}
	span := &Span{
		name:    name,
		Level:   level,
		started: time.Now(),
		Ctx:     ctx,
		Num:     num,
	}

	return context.WithValue(ctx, "span", span), span
}

func (s *Span) End() {
	s.ended = common.Time(time.Now())
	if s.status == "" {
		s.status = Ok
	}
}

func (s *Span) SetStatus(code StatusCode, description string) {
	s.status = code
	s.description = description
}

func (s *Span) TimeEstimate() time.Duration {
	if s.ended != nil {
		return s.ended.Sub(s.started)
	}
	return 0
}

func (s *Span) SetAttributes(key string, value interface{}) {
	if s.attributes == nil {
		s.attributes = make(map[string]string)
	}
	s.attributes[key] = fmt.Sprintf("%v", value)
}

func SpanFromContext(ctx context.Context) (span *Span, ok bool) {
	if ctx == nil {
		return
	}
	span, ok = ctx.Value("span").(*Span)
	return
}

type StateItem struct {
	Name         string
	Num          int
	Start        time.Time
	End          *time.Time
	TimeEstimate time.Duration
	Attributes   map[string]string
	Status       StatusCode
	Level        int
}

type Telemetry []*StateItem

func (s Telemetry) Len() int           { return len(s) }
func (s Telemetry) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Telemetry) Less(i, j int) bool { return s[i].Num < s[j].Num }

func Unpack(ctx context.Context) Telemetry {
	_levels := make(Telemetry, 0)
LOOP:
	span, ok := SpanFromContext(ctx)
	if !ok {
		return _levels
	}

	_levels = append(_levels, &StateItem{
		Name:         span.name,
		Num:          span.Num,
		Start:        span.started,
		End:          span.ended,
		TimeEstimate: span.TimeEstimate(),
		Attributes:   span.attributes,
		Status:       span.status,
		Level:        span.Level,
	})

	sort.Sort(_levels)

	if span.Ctx != nil {
		ctx = span.Ctx
		goto LOOP
	}

	return nil
}
