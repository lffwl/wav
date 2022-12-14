package wav

import (
	"fmt"
	"io"
	"strings"
)

type riff struct {
	io.ReadWriteSeeker
	format Format
	other  map[string][]byte
}

func newRead(seeker io.ReadWriteSeeker) Wav {
	r := &riff{
		ReadWriteSeeker: seeker,
		other:           map[string][]byte{},
	}

	// parse wav head
	r.parse()

	return r
}

// Format output format
func (r *riff) Format() *Format {
	return &r.format
}

// readSizeByte
// size length
func (r *riff) readSizeByte(size uint32) []byte {
	data := make([]byte, size)
	if n, err := r.Read(data); err != nil {
		panic(fmt.Sprintf("wav read failed : %v", err))
	} else if n != int(size) {
		panic("the read data length is wrong")
	}
	return data
}

func (r *riff) parse() {
	// discard the first 12 bytes
	// discard RIFF chunk
	if _, err := r.Seek(12, io.SeekStart); err != nil {
		panic(fmt.Sprintf("seek failed : %v", err))
	}

	// parse format chunk
	r.parseFormat()

	// parse other chunk
	r.parseOther()
}

// parseFormat parse format chunk
func (r *riff) parseFormat() {
	r.format.id = string(r.readSizeByte(4)[:])
	r.format.Size = BytesToUint32(r.readSizeByte(4))
	r.format.AudioFormat = BytesToUint16(r.readSizeByte(2))
	r.format.NumChannels = BytesToUint16(r.readSizeByte(2))
	r.format.SampleRate = BytesToUint32(r.readSizeByte(4))
	r.format.ByteRate = BytesToUint32(r.readSizeByte(4))
	r.format.BlockAlign = BytesToUint16(r.readSizeByte(2))
	r.format.BitsPerSample = BytesToUint16(r.readSizeByte(2))
}

// parseOther parse other chunk
func (r *riff) parseOther() {
	for {
		chunkId := string(r.readSizeByte(4)[:])
		// Read the data chunk id exit for
		if chunkId == DataChunkId {
			return
		} else {
			size := BytesToUint32(r.readSizeByte(4))
			// Prevent chunk id length less than 4
			// fill in the blanks
			r.other[strings.TrimSpace(chunkId)] = r.readSizeByte(size)
		}
	}
}
