package chunk

import (
	"encoding/binary"
	"errors"
	"io"
)

// RiffChunk is a structure that handles riff chunk descriptor of wave.
type RiffChunk struct {
	ID     string
	Size   uint32
	Format string
}

func NewRiffChunk(file io.Reader) (*RiffChunk, error) {
	const chunkByteSize = 12
	chunkBytes := make([]byte, chunkByteSize)
	if _, err := io.ReadFull(file, chunkBytes); err != nil {
		return nil, err
	}

	chunk := &RiffChunk{
		ID:     string(chunkBytes[:4]),
		Size:   binary.LittleEndian.Uint32(chunkBytes[4:]),
		Format: string(chunkBytes[8:]),
	}
	if err := chunk.validate(); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (chunk *RiffChunk) validate() error {
	if chunk.ID != "RIFF" {
		return errors.New("RiffChunk: ChunkID must be [RIFF]")
	}
	if chunk.Format != "WAVE" {
		return errors.New("RiffChunk: Format must be [WAVE]")
	}
	return nil
}
