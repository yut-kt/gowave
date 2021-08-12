package chunk

import (
	"encoding/binary"
	"errors"
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

	chunk := &RiffChunk{
		id:     string(chunkBytes[:4]),
		size:   binary.LittleEndian.Uint32(chunkBytes[4:]),
		format: string(chunkBytes[8:]),
	}
	if err := chunk.validate(); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (chunk *RiffChunk) validate() error {
	if chunk.id != "RIFF" {
		return errors.New("RiffChunk: ChunkID must be [RIFF]")
	}
	if chunk.format != "WAVE" {
		return errors.New("RiffChunk: Format must be [WAVE]")
	}
	return nil
}
