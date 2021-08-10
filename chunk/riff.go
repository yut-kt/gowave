package chunk

import (
	"encoding/binary"
	"io"
)

type RiffChunk struct {
	id     string
	size   uint32
	format string
}

func NewRiffChunk(file io.Reader) (*RiffChunk, error) {
	const chunkByteSize = 12
	chunkBytes := make([]byte, chunkByteSize)
	if _, err := io.ReadFull(file, chunkBytes); err != nil {
		return nil, err
	}

	return &RiffChunk{
		id:     string(chunkBytes[:4]),
		size:   binary.LittleEndian.Uint32(chunkBytes[4:]),
		format: string(chunkBytes[8:]),
	}, nil
}
