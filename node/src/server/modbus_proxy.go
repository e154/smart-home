package server

import (
	"../cache"
	"../serial"
	"fmt"
)

const (
	ADDRESS uint8 = 0
)

func ModBusProxy(conn *serial.Serial, command []byte) (res []byte, err error) {

	cache_ptr := cache.CachePtr()
	cache_key := cache_ptr.GetKey(fmt.Sprintf("%d_dev", command[ADDRESS]))


	for i := 0; i<5; i++ {

		cache_exist := cache_ptr.IsExist(cache_key)
		if cache_exist {
			conn.Dev = cache_ptr.Get(cache_key).(string)
			res, err = ModBusExec(conn, command)
			if err == nil {
				return
			}
		} else {

			devices := serial.FindSerials()
			for _, device := range devices {
				conn.Dev = device
				res, err = ModBusExec(conn, command)
				if err == nil {
					return
				}
			}
		}

	}

	return
}

func ModBusExec(conn *serial.Serial, command []byte) (res []byte, err error) {

	// get cache
	cache_ptr := cache.CachePtr()
	cache_key := cache_ptr.GetKey(fmt.Sprintf("%d_dev", command[ADDRESS]))

	if _, err = conn.Open(); err != nil {
		cache_ptr.Delete(cache_key)
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}

	modbus := &serial.Modbus{Serial: conn}
	res, err = modbus.Send(command)
	if err != nil {
		cache_ptr.Delete(cache_key)
		//log.Printf("error: %s - %s\r\n",conn.Dev, err.Error())
		return
	}

	//cache update
	cache_ptr.Put("node", cache_key, conn.Dev)

	return
}