package rpc

import "time"

const (
	ADDRESS 	uint8 = iota
	COMMAND
	DATA
)

type Request struct {
	Line		string		`json: "line"`
	Device		string		`json: "device"`
	Baud		int		`json: "baud"`
	StopBits	int		`json: "stop_bits"`
	Sleep		int		`json: "sleep"`
	Timeout		time.Duration	`json: "timeout"`
	Command		[]byte		`json: "command"`
	Result		bool		`json: "result"`
}

type Result struct {
	Command   []byte		`json: "command"`
	Device    string		`json: "device"`
	Result    string		`json: "result"`
	Error     string		`json: "error"`
	ErrorCode string		`json: "error_code"`
}