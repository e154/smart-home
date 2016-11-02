package pack

import "time"

const (
	ADDRESS 	uint8 = iota
	COMMAND
	DATA
)

type Package struct {
	Line		string		`json: "line"`
	Device		string		`json: "device"`
	Baud		int		`json: "baud"`
	StopBits	int		`json: "stop_bits"`
	Timeout		time.Duration	`json: "timeout"`
	Command		[]byte		`json: "command"`
	Result		bool		`json: "result"`
	Time		time.Time	`json: "time"`
}

type Result struct {
	Command		[]byte		`json: "command"`
	Line		string		`json: "line"`
	Device		string		`json: "device"`
	Result		[]byte		`json: "result"`
	Time		time.Time	`json: "time"`
	Error		string		`json: "error"`
	ErrorCode	int		`json: "error_code"`
}