package wav

import (
	"fmt"
	"io"
	"strings"
)

type riff struct {
	io.ReadWriteSeeker
	format  Format
	dataLen uint32
	other   map[string][]byte
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

// newWrite
// seeker io.ReadWriteSeeker
// dataLen data length 音频数据的长度
// numChannels 音频数据的声道数，1：单声道，2：双声道
// sampleRate 音频数据的采样率
// bitsPerSample 每个采样存储的bit数
// audioFormat Data区块存储的音频数据的格式，PCM音频数据的值为1
// other chunk
func newWrite(seeker io.ReadWriteSeeker, dataLen uint32, numChannels uint16, sampleRate uint32, bitsPerSample uint16, audioFormat uint16, other ...map[string][]byte) Wav {
	r := &riff{
		ReadWriteSeeker: seeker,
		dataLen:         dataLen,
	}

	// set format
	r.format = Format{
		id:            FmtChunkId,
		Size:          16,
		AudioFormat:   audioFormat,
		NumChannels:   numChannels,
		SampleRate:    sampleRate,
		ByteRate:      uint32(int(sampleRate)*int(numChannels)*int(bitsPerSample)) / 8,
		BlockAlign:    uint16(int(numChannels)*int(bitsPerSample)) / 8,
		BitsPerSample: bitsPerSample,
	}

	// set other
	if len(other) > 0 {
		r.other = other[0]
	}

	// write head
	r.writeHead()

	return r
}

// Format output format
func (r *riff) Format() *Format {
	return &r.format
}

// GetDataLen output data length
func (r *riff) GetDataLen() uint32 {
	return r.dataLen
}

// writeHead
// write wav head
func (r *riff) writeHead() {
	r.writeRIFFChunk()
	r.writeFormatChunk()
	r.writeOtherChunk()
	// write data head
	if n, err := r.Write([]byte(DataChunkId)); err != nil || n != 4 {
		panic("writer wav data chunk id failed")
	}
	if n, err := r.Write(Uint32ToBytes(r.dataLen)); err != nil || n != 4 {
		panic("writer wav data chunk size failed")
	}
}

// writeRIFFChunk
// write RIFF Chunk
func (r *riff) writeRIFFChunk() {
	if n, err := r.Write([]byte(RIFFChunkId)); err != nil || n != 4 {
		panic("writer wav RIFF chunk id failed")
	}
	// RIFF Size
	// data length + head length - RIFF ID length - RIFF size length + other length
	size := r.dataLen + 44 - 8 + r.getOtherLen()
	if n, err := r.Write(Uint32ToBytes(size)); err != nil || n != 4 {
		panic("writer wav RIFF chunk size failed")
	}
	if n, err := r.Write([]byte(RIFFChunkType)); err != nil || n != 4 {
		panic("writer wav  RIFF chunk type failed")
	}
}

// writeFormatChunk
// write Format Chunk
func (r *riff) writeFormatChunk() {
	if n, err := r.Write([]byte(FmtChunkId)); err != nil || n != 4 {
		panic("writer wav Format chunk id failed")
	}
	if n, err := r.Write(Uint32ToBytes(16)); err != nil || n != 4 {
		panic("writer wav Format chunk size failed")
	}
	if n, err := r.Write(Uint16ToBytes(r.format.AudioFormat)); err != nil || n != 2 {
		panic("writer wav Format chunk AudioFormat failed")
	}
	if n, err := r.Write(Uint16ToBytes(r.format.NumChannels)); err != nil || n != 2 {
		panic("writer wav Format chunk NumChannels failed")
	}
	if n, err := r.Write(Uint32ToBytes(r.format.SampleRate)); err != nil || n != 4 {
		panic("writer wav Format chunk SampleRate failed")
	}
	if n, err := r.Write(Uint32ToBytes(r.format.ByteRate)); err != nil || n != 4 {
		panic("writer wav Format chunk ByteRate failed")
	}
	if n, err := r.Write(Uint16ToBytes(r.format.BlockAlign)); err != nil || n != 2 {
		panic("writer wav Format chunk BlockAlign failed")
	}
	if n, err := r.Write(Uint16ToBytes(r.format.BitsPerSample)); err != nil || n != 2 {
		panic("writer wav Format chunk BitsPerSample failed")
	}
}

// writeFormatChunk
// write Other Chunk
func (r *riff) writeOtherChunk() {
	if len(r.other) > 0 {
		for name, bytes := range r.other {
			if n, err := r.Write(OtherNameToBytes(name)); err != nil || n != 4 {
				panic(fmt.Sprintf("writer wav Other chunk Name %s failed", name))
			}
			if n, err := r.Write(Uint32ToBytes(uint32(len(bytes)))); err != nil || n != 4 {
				panic("writer wav Other chunk size failed")
			}
			if n, err := r.Write(bytes); err != nil || n != len(bytes) {
				panic("writer wav Other chunk bytes failed")
			}
		}
	}
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
			r.dataLen = BytesToUint32(r.readSizeByte(4))
			return
		} else {
			size := BytesToUint32(r.readSizeByte(4))
			// Prevent chunk id length less than 4
			// fill in the blanks
			r.other[strings.TrimSpace(chunkId)] = r.readSizeByte(size)
		}
	}
}

// getOtherLen
// get other bytes length
func (r *riff) getOtherLen() uint32 {
	if len(r.other) > 0 {
		var size = 0
		for _, bytes := range r.other {
			// 8 len = name 4 byte + size 4 byte
			// + len(bytes)
			size += 8 + len(bytes)
		}
		return uint32(size)
	}
	return 0
}
