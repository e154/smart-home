package serial

import (
	"bufio"
	"fmt"
	"bytes"
)

const (
	STATE_RX_IDLE uint8 = iota
	STATE_RX_RCV
	STATE_RX_WAIT_EOF
)

const (
	BYTE_HIGH_NIBBLE uint8 = iota
	BYTE_LOW_NIBBLE
)

const (
	hexTable = "0123456789ABCDEF"
)

type Modbus struct {
	serial *Serial
	rcvState uint8
	rcvBytePos uint8
	rcvBufferPos uint8
	rcvBuf []byte
	trcBuff *bytes.Buffer
}

// 1 - address		u08
// 1 - function		u08
// 1..n - data		u08 x N
func (m *Modbus) Send(data []byte) (result []byte, err error) {

	lrc := LRC(data)
	data = append(data, lrc)

	if err = m.asciiTransmit(data); err != nil {
		return
	}

	m.rcvState = STATE_RX_IDLE
	reader := bufio.NewReader(m.serial.Port)
	for {
		b, err := reader.ReadByte();
		if err != nil {
			return result, err
			break
		}

		result, err = m.asciiReceiveFSM(b)
		if err != nil {
			return result, err
			break
		}

		if m.rcvState == STATE_RX_IDLE {
			result = m.rcvBuf
			break
		}

	}

	return
}

func (m *Modbus) asciiReceiveFSM(b byte) ([]byte, error) {

	fmt.Print(string(b))

	switch m.rcvState {
	case STATE_RX_RCV:
		if( b == ':' ) {
			m.rcvBuf = []byte{}
			m.rcvBytePos = BYTE_HIGH_NIBBLE;
		} else if( b == '\r' ) {
			m.rcvState = STATE_RX_WAIT_EOF;
		} else {
			switch m.rcvBytePos {
			case BYTE_HIGH_NIBBLE:
				m.rcvBuf = append(m.rcvBuf, b<<4)
				m.rcvBytePos = BYTE_LOW_NIBBLE
			case BYTE_LOW_NIBBLE:
				m.rcvBuf[len(m.rcvBuf) - 1] |= b
				m.rcvBytePos = BYTE_HIGH_NIBBLE
			}
		}
	case STATE_RX_WAIT_EOF:
		if (b == '\n') {
			m.rcvState = STATE_RX_IDLE
			return m.rcvBuf, nil
		} else if (b == ':') {
			m.rcvBytePos = BYTE_HIGH_NIBBLE;
			m.rcvState = STATE_RX_RCV;
		} else {
			m.rcvState = STATE_RX_IDLE;
		}
	case STATE_RX_IDLE:
		if (b == ':') {
			m.rcvBuf = []byte{}
			m.rcvBytePos = BYTE_HIGH_NIBBLE;
			m.rcvState = STATE_RX_RCV;
		}
	}

	return  []byte{}, nil
}

// 1 - address		u08
// 1 - function		u08
// 1..n - data		u08 x N
// 1 - lrc		u08
// 1 - \r		u08
// 1 - \n		u08
func (m *Modbus) asciiTransmit(data []byte) error {

//TODO add decoder
	if m.trcBuff != nil {
		m.trcBuff.Reset()
	}

	m.trcBuff = &bytes.Buffer{}
	m.trcBuff.WriteByte(':')

	for _, d := range data {
		m.trcBuff.WriteByte(HI(d))
		m.trcBuff.WriteByte(LOW(d))
	}

	m.trcBuff.Write([]byte{'\r', '\n'})

	if m.serial.Port == nil {
		if _, err := m.serial.Open(); err != nil {
			return err
		}
	}

	//fmt.Printf("%X\r\n", buf)

	_, err := m.serial.Port.Write(m.trcBuff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (m *Modbus) Device(address byte) *Device {
	return &Device{address: address, modbus: m}
}

func (m *Modbus) bin2char(buf *bytes.Buffer, b byte) error {

	var h [2]byte
	h[0] = hexTable[b>>4]
	h[1] = hexTable[b&0xf]

	if _, err := buf.Write(h[:]); err != nil {
		return err
	}

	return nil
}

func (m *Modbus) char2bin(b byte) byte {

	return b
}

func (m *Modbus) error() {

}