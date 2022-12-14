package wav

import "io"

type Wav interface {
	io.ReadWriteSeeker
	Format() *Format
}

type Format struct {
	id            string
	Size          uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}
