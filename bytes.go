package wav

import "encoding/binary"

func BytesToUint32(v []byte) uint32 {
	return binary.LittleEndian.Uint32(v)
}

func BytesToUint16(v []byte) uint16 {
	return binary.LittleEndian.Uint16(v)
}
