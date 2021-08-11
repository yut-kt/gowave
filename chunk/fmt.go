package chunk

import (
	"encoding/binary"
	"errors"
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
	chunk := &FmtChunk{
		id:            string(chunkBytes[:4]),
		size:          binary.LittleEndian.Uint32(chunkBytes[4:]),
		audioFormat:   binary.LittleEndian.Uint16(chunkBytes[8:]),
		numChannels:   binary.LittleEndian.Uint16(chunkBytes[10:]),
		sampleRate:    binary.LittleEndian.Uint32(chunkBytes[12:]),
		byteRate:      binary.LittleEndian.Uint32(chunkBytes[16:]),
		blockAlign:    binary.LittleEndian.Uint16(chunkBytes[20:]),
		bitsPerSample: binary.LittleEndian.Uint16(chunkBytes[22:]),
	}
	if err := chunk.validate(); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (chunk *FmtChunk) validate() error {
	if chunk.id != "fmt " {
		return errors.New("FmtChunk: SubChunk1ID must be [fmt ]")
	}
	if chunk.size != 16 {
		return errors.New("FmtChunk: SubChunk1Size still only supports 16 bytes")
	}
	if chunk.audioFormat != 1 {
		return errors.New("FmtChunk: AudioFormat still only supports 1(PCM)")
	}
	if chunk.numChannels != 1 {
		return errors.New("FmtChunk: NumChannels still only supports 1(Mono)")
	}
	if chunk.bitsPerSample != 8 && chunk.bitsPerSample != 16 {
		return errors.New("FmtChunk: BitsPerSample still only supports 8 or 16")
	}
	return nil
}

func (chunk *FmtChunk) GetBitsPerSample() uint16 {
	return chunk.bitsPerSample
}
