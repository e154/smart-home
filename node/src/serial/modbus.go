package serial

import (
	"bufio"
	"fmt"
	"bytes"
	"errors"
)

const (
	ILLEGAL_FUNCTION uint8 = iota + 1
	ILLEGAL_DATA_ADDRESS
	ILLEGAL_DATA_VALUE
	SLAVE_DEVICE_FAILURE
	ACKNOWLEDGE
	SLAVE_DEVICE_BUSY
	NEGATIVE_ACKNOWLEDGE
	MEMORY_PARITY_ERROR
	ILLEGAL_LRC
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

	//fmt.Print(string(b))
	switch m.rcvState {
	case STATE_RX_RCV:
		if( b == ':' ) {
			m.rcvBuf = []byte{}
			m.rcvBytePos = BYTE_HIGH_NIBBLE;
		} else if( b == '\r' ) {
			m.rcvState = STATE_RX_WAIT_EOF;
		} else {
			b = m.char2bin(b)
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
			fmt.Printf("receive <- %X\r\n", m.rcvBuf) //TODO remove

			return m.rcvBuf, m.checkError(m.rcvBuf)
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

	fmt.Printf("send -> %X\r\n", m.trcBuff.Bytes()) //TODO remove

	_, err := m.serial.Port.Write(m.trcBuff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (m *Modbus) Device(address byte) *Device {
	return &Device{address: address, modbus: m}
}

func (m *Modbus) bin2char(b byte) byte {

	if( b <= 0x09 ) {
		return byte( '0' + b )
	} else if( ( b >= 0x0A ) && ( b <= 0x0F ) ) {
		return byte( b - 0x0A + 'A' )
	}

	return '0'
}

func (m *Modbus) char2bin(b byte) byte {

	if( b >= '0' ) && ( b <= '9' ) {
		return byte( b - '0' )
	} else if( ( b >= 'A' ) && ( b <= 'F' ) ) {
		return byte( b - 'A' + 0x0A )
	}

	return 0xFF
}

func (m *Modbus) checkError(buf []byte) error {

	var errCode uint8

	// check lrc
	if ( LRC(buf[0:len(buf) - 1]) != buf[len(buf) - 1] ) {
		errCode = ILLEGAL_LRC
	}

	// check error bite
	if ( buf[1] & (1<<7) != 0 ) {
		errCode = buf[2]
	}

	// convert error code
	var err error
	switch errCode {
	case ILLEGAL_FUNCTION:
		err = errors.New("ILLEGAL_FUNCTION")
	case ILLEGAL_DATA_ADDRESS:
		err = errors.New("ILLEGAL_DATA_ADDRESS")
	case ILLEGAL_DATA_VALUE:
		err = errors.New("ILLEGAL_DATA_VALUE")
	case SLAVE_DEVICE_FAILURE:
		err = errors.New("SLAVE_DEVICE_FAILURE")
	case ACKNOWLEDGE:
		err = errors.New("ACKNOWLEDGE")
	case SLAVE_DEVICE_BUSY:
		err = errors.New("SLAVE_DEVICE_BUSY")
	case NEGATIVE_ACKNOWLEDGE:
		err = errors.New("NEGATIVE_ACKNOWLEDGE")
	case MEMORY_PARITY_ERROR:
		err = errors.New("MEMORY_PARITY_ERROR")
	case ILLEGAL_LRC:
		err = errors.New("ILLEGAL_LRC")
	}

	return err
}