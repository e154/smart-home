package serial

import (
	"github.com/tarm/serial"
	"time"
)

type Serial struct {
	Dev		string
	Baud		int
	ReadTimeout 	time.Duration
	StopBits 	int
	config		*serial.Config
	Port		*serial.Port
}

func (s *Serial) Open() (*Serial, error) {

	s.config = &serial.Config{
		Name: s.Dev,
		Baud: s.Baud,
		StopBits: serial.StopBits(s.StopBits),
		ReadTimeout: s.ReadTimeout,
	}

	var err error
	if s.Port, err = serial.OpenPort(s.config); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Serial) Close() (*Serial, error) {

	if s.Port != nil {
		return s, s.Port.Close()
	}

	return s, nil
}