package wav

import "encoding/binary"

func BytesToUint32(v []byte) uint32 {
	return binary.LittleEndian.Uint32(v)
}

func BytesToUint16(v []byte) uint16 {
	return binary.LittleEndian.Uint16(v)
}

func Uint32ToBytes(v uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, v)
	return bytes
}

func Uint16ToBytes(v uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, v)
	return bytes
}

func OtherNameToBytes(v string) []byte {
	if len(v) < 4 {
		for i := 0; i <= 4-len(v); i++ {
			v += " "
		}
	} else if len(v) > 4 {
		v = v[0:4]
	}
	bytes := []byte(v)
	return bytes
}
