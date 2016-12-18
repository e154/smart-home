package serial

import (
	"bufio"
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
	Serial       *Serial
	rcvState     uint8
	rcvBytePos   uint8
	rcvBufferPos uint8
	rcvBuf       []byte
	trcBuff      *bytes.Buffer
}

// 1 - address		u08
// 1 - function		u08
// 1..n - data		u08 x N
func (m *Modbus) Send(data []byte) (result []byte, err error) {

	var b byte
	var ok bool
	lrc := LRC(data)
	data = append(data, lrc)

	if err = m.asciiTransmit(data); err != nil {
		return
	}

	m.rcvState = STATE_RX_IDLE
	reader := bufio.NewReader(m.Serial.Port)
	for {
		if b, err = reader.ReadByte(); err != nil {
			break
		}

		if ok, err = m.asciiReceiveFSM(b); err != nil || ok {
			break
		}
	}

	if len(m.rcvBuf) > 2 {
		result = m.rcvBuf[2:len(m.rcvBuf) - 1]
	}

	//log.Printf("receive <- %X, len %d\r\n", result, len(m.rcvBuf)) //TODO remove

	return
}

func (m *Modbus) asciiReceiveFSM(b byte) (bool, error) {

	//log.Print(string(b))
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
			//log.Printf("receive <- %X, len: %d\r\n", m.rcvBuf, len(m.rcvBuf)) //TODO remove
			return true, m.checkError(m.rcvBuf)

		} else if (b == ':') {
			m.rcvBuf = []byte{}
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

	return  false, nil
}

// 1 - address		u08
// 1 - function		u08
// 1..n - data		u08 x N
// 1 - lrc		u08
// 1 - \r		u08
// 1 - \n		u08
func (m *Modbus) asciiTransmit(data []byte) (err error) {

	if m.trcBuff != nil {
		m.trcBuff.Reset()
	}

	m.trcBuff = &bytes.Buffer{}
	if err = m.trcBuff.WriteByte(':'); err != nil {
		return
	}

	for _, d := range data {
		m.trcBuff.WriteByte(m.bin2char(HI(d)))
		m.trcBuff.WriteByte(m.bin2char(LOW(d)))
	}

	m.trcBuff.Write([]byte{'\r', '\n'})

	if m.Serial == nil {
		err = errors.New("serial pointer is nil")
		return
	}

	//log.Printf("send -> %X\r\n", m.trcBuff.Bytes()) //TODO comment

	_, err = m.Serial.Port.Write(m.trcBuff.Bytes())
	if err != nil {
		return
	}

	return
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