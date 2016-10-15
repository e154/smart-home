package serial

import (
	"log"
	"bufio"
	"os"

	"github.com/tarm/serial"
	"fmt"
	"io/ioutil"
	"strings"
	"io"
)

const (
	BAUD int = 9600
	PORT string = "/dev/ttyUSB0"
)

var (
	PORTS []string
)

func checkErr(err error) {

	if err != nil {
		log.Fatal(err.Error)
		os.Exit(1)
	}
}

func getBytes(buf *bufio.Reader, n int) []byte {
	// Читаем n байт
	bytes, err:= buf.Peek(n)
	checkErr(err)
	// Освобождаем n байт
	//skipBytes(buf, n)
	return bytes
}



func Init() {
	c := &serial.Config{Name: PORT, Baud: BAUD, StopBits: serial.Stop2}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	//n, err := s.Write([]byte("test"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	findSerials()

	fmt.Println(PORTS)

	r := bufio.NewReader(s)
	//for {
		buf := getBytes(r, 16)
		fmt.Println(buf)
	//}

}

func findSerials() {
	contents, _ := ioutil.ReadDir("/dev")

	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyS") ||
			strings.Contains(f.Name(), "ttyUSB") {
			PORTS = append(PORTS, "/dev/" + f.Name())
		}
	}

}

func sendCommand(command byte, argument float32, serialPort io.ReadWriteCloser) error {


	return nil
}