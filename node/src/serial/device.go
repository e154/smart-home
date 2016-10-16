package serial

type Device struct {
	address	byte
	modbus	*Modbus
}

func (m *Device) Send(data []byte) (result []byte, err error) {

	var ndata []byte

	for i, d := range data {
		ndata[i+1] = d
	}

	ndata[0] = m.address

	return m.modbus.Send(ndata)
}