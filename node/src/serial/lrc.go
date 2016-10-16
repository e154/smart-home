package serial

func LRC(data []byte) byte {

	var ucLRC uint8 = 0

	var b byte
	for _, b = range data {
		ucLRC += b
	}

	return uint8(-int8(ucLRC))
}
