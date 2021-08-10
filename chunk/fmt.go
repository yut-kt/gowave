package chunk

import (
	"encoding/binary"
	"io"
)

type FmtChunk struct {
	id            string
	size          uint32
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16
}

func NewFmtChunk(file io.Reader) (*FmtChunk, error) {
	const chunkByteSize = 24
	chunkBytes := make([]byte, chunkByteSize)
	if _, err := io.ReadFull(file, chunkBytes); err != nil {
		return nil, err
	}
	return &FmtChunk{
		id:            string(chunkBytes[:4]),
		size:          binary.LittleEndian.Uint32(chunkBytes[4:]),
		audioFormat:   binary.LittleEndian.Uint16(chunkBytes[8:]),
		numChannels:   binary.LittleEndian.Uint16(chunkBytes[10:]),
		sampleRate:    binary.LittleEndian.Uint32(chunkBytes[12:]),
		byteRate:      binary.LittleEndian.Uint32(chunkBytes[16:]),
		blockAlign:    binary.LittleEndian.Uint16(chunkBytes[20:]),
		bitsPerSample: binary.LittleEndian.Uint16(chunkBytes[22:]),
	}, nil
}

func (fmtChunk *FmtChunk) GetBitsPerSample() uint16 {
	return fmtChunk.bitsPerSample
}
