package wav

import (
	"errors"
	"io"
)

func NewRead(seeker io.ReadWriteSeeker) (Wav, error) {
	chunkId := make([]byte, 4)
	if n, err := seeker.Read(chunkId); err != nil || n != 4 {
		return nil, errors.New("wav read chunk id failed")
	}

	var w Wav
	switch string(chunkId[:]) {
	case RIFFChunkId:
		w = newRead(seeker)
	default:
		return nil, errors.New("this type is not supported yet")
	}

	return w, nil
}
