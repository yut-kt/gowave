package chunk

import (
	"encoding/binary"
	"errors"
	"io"
)

// FmtChunk is a structure that handles fmt subChunk of wave.
type FmtChunk struct {
	ID            string
	Size          uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

// NewFmtChunk is a function to construct FmtChunk struct.
func NewFmtChunk(file io.Reader) (*FmtChunk, error) {
	const chunkByteSize = 24
	chunkBytes := make([]byte, chunkByteSize)
	if _, err := io.ReadFull(file, chunkBytes); err != nil {
		return nil, err
	}
	chunk := &FmtChunk{
		ID:            string(chunkBytes[:4]),
		Size:          binary.LittleEndian.Uint32(chunkBytes[4:]),
		AudioFormat:   binary.LittleEndian.Uint16(chunkBytes[8:]),
		NumChannels:   binary.LittleEndian.Uint16(chunkBytes[10:]),
		SampleRate:    binary.LittleEndian.Uint32(chunkBytes[12:]),
		ByteRate:      binary.LittleEndian.Uint32(chunkBytes[16:]),
		BlockAlign:    binary.LittleEndian.Uint16(chunkBytes[20:]),
		BitsPerSample: binary.LittleEndian.Uint16(chunkBytes[22:]),
	}
	if err := chunk.validate(); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (chunk *FmtChunk) validate() error {
	if chunk.ID != "fmt " {
		return errors.New("FmtChunk: SubChunk1ID must be [fmt ]")
	}
	if chunk.Size != 16 {
		return errors.New("FmtChunk: SubChunk1Size still only supports 16 bytes")
	}
	if chunk.AudioFormat != 1 {
		return errors.New("FmtChunk: AudioFormat still only supports 1(PCM)")
	}
	if chunk.NumChannels != 1 {
		return errors.New("FmtChunk: NumChannels still only supports 1(Mono)")
	}
	if chunk.BitsPerSample != 8 && chunk.BitsPerSample != 16 {
		return errors.New("FmtChunk: BitsPerSample still only supports 8 or 16")
	}
	return nil
}
