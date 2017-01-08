package server

import (
	"../cache"
	"../serial"
	"../settings"
	"../lib/rpc"
	"fmt"
	"errors"
	"encoding/hex"
	"sync"
	"time"
)

const (
	ADDRESS uint8 = 0
)

type Modbus struct {
	mu	sync.Mutex
}

func (m *Modbus) Send(request *rpc.Request, result *rpc.Result) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(request.Command) == 0 {
		return errors.New("command == []")
	}

	st := settings.SettingsPtr()

	conn := &serial.Serial{
		Dev: "",
		Baud: st.Baud,
		ReadTimeout: st.Timeout,
		StopBits: st.StopBits,
	}

	if request.Device != "" {
		conn.Dev = request.Device
	}

	if request.Baud != 0 {
		conn.Baud = request.Baud
	}

	if request.Timeout != 0 {
		conn.ReadTimeout = request.Timeout
	}

	var err error

	if conn.Dev == "" {

		cache_ptr := cache.CachePtr()
		cache_key := cache_ptr.GetKey(fmt.Sprintf("%d_dev", request.Command[ADDRESS]))

		//log.Println("send", request.Command)
		//for i := 0; i<5; i++ {

			cache_exist := cache_ptr.IsExist(cache_key)
			if cache_exist {
				conn.Dev = cache_ptr.Get(cache_key).(string)
				result.Result, err, result.ErrorCode = m.exec(conn, request)
				if err == nil {
					result.Device = conn.Dev
					return nil
				}
			} else {

				//
				devices := serial.FindSerials()
				for _, device := range devices {
					conn.Dev = device
					result.Result, err, result.ErrorCode = m.exec(conn, request)
					if err == nil {
						result.Device = device
						return nil
					}
				}
			}

		//}
	} else {
		for i := 0; i<5; i++ {
			result.Result, err, result.ErrorCode = m.exec(conn, request)
			if err == nil {
				result.Device = conn.Dev
				return nil
			}
		}
	}

	if err != nil {
		result.Error = err.Error()
	}

	return nil
}

func (m *Modbus) exec(conn *serial.Serial, request *rpc.Request) (result string, err error, errcode string) {

	// get cache
	cache_ptr := cache.CachePtr()
	cache_key := cache_ptr.GetKey(fmt.Sprintf("%d_dev", request.Command[ADDRESS]))

	if _, err = conn.Open(); err != nil {
		//cache_ptr.Delete(cache_key)
		errcode = "SERIAL_PORT_ERROR"
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}
	defer conn.Close()

	modbus := &serial.Modbus{Serial: conn}
	var b []byte
	if b, err = modbus.Send(request.Command); err != nil {
		//cache_ptr.Delete(cache_key)
		errcode = "MODBUS_LINE_ERROR"
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}

	result = hex.EncodeToString(b)
	cache_ptr.Put("node", cache_key, conn.Dev)

	// bug in the devices need timeout, need fix!!!
	if request.Sleep != 0 {
		time.Sleep(time.Millisecond * time.Duration(request.Sleep))
	}

	return
}