package web

import (
	"time"
)

// Request ...
type Request struct {
	Method  string
	Url     string
	Body    []byte
	Headers []map[string]string
	Timeout time.Duration
}

type Crawler interface {
	Probe(Request) (int, []byte, error)
}
