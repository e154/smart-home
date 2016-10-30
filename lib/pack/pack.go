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
	Timeout		time.Duration	`json: "timeout"`
	Command		[]byte		`json: "command"`
	Result		bool		`json: "result"`
	Time		time.Time	`json: "time"`
}

type Error struct {
	Error		string		`json :"error"`
	Time		time.Time	`json: "time"`
}

type Result struct {
	Command		[]byte		`json: "command"`
	Line		string		`json: "line"`
	Device		string		`json: "device"`
	Result		[]byte		`json: "result"`
	Time		time.Time	`json: "time"`
}