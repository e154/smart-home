package serial
//
//import (
//	"github.com/goburrow/serial"
//)
//
//type Serial2 struct {
//	config serial.Config
//	Port serial.Port
//	isConnected bool
//}
//
//func (s *Serial2) Connect() (err error) {
//	if s.isConnected {
//		return
//	}
//	if s.Port == nil {
//		s.Port, err = serial.Open(&s.config)
//	} else {
//		err = s.Port.Open(&s.config)
//	}
//	if err == nil {
//		s.isConnected = true
//	}
//	return
//}
//
//func (s *Serial2) Close() (err error) {
//	if !s.isConnected {
//		return
//	}
//	err = s.Port.Close()
//	s.isConnected = false
//	return
//}
//