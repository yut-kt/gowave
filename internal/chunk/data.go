package chunk

import (
	"encoding/binary"
	"errors"
	"io"
)

// DataChunk is a structure that handles data subChunk of wave.
type DataChunk struct {
	File io.ReadSeeker
	ID   string
	Size uint32
	Data interface{}
}

// NewDataChunk is a function to construct DataChunk struct.
func NewDataChunk(file io.ReadSeeker) (*DataChunk, error) {
	const chunkHeaderByteSize = 8
	chunkHeaderBytes := make([]byte, chunkHeaderByteSize)
	if _, err := io.ReadFull(file, chunkHeaderBytes); err != nil {
		return nil, err
	}

	chunk := &DataChunk{
		File: file,
		ID:   string(chunkHeaderBytes[:4]),
		Size: binary.LittleEndian.Uint32(chunkHeaderBytes[4:]),
	}
	if err := chunk.validate(); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (chunk *DataChunk) validate() error {
	if chunk.ID != "data" {
		return errors.New("DataChunk: ChunkID must be [data]")
	}
	return nil
}

// ReadData is a function to read the sample in the wave.
func (chunk *DataChunk) ReadData(bitsPerSample uint16, samplingNum int64) (data interface{}, err error) {
	switch bitsPerSample {
	case 8:
		data, err = readSamples[uint8](chunk.File, samplingNum, 1)
	case 16:
		data, err = readSamples[int16](chunk.File, samplingNum, 2)
	}
	if err != nil {
		return
	}

	switch v := chunk.Data.(type) {
	case []uint8:
		d, _ := data.([]uint8)
		chunk.Data = append(v, d...)
	case []int16:
		d, _ := data.([]int16)
		chunk.Data = append(v, d...)
	default:
		chunk.Data = data
	}
	return
}

type Sample interface {
	uint8 | int16
}

func readSamples[S Sample](f io.ReadSeeker, samplingN, samplingBytes int64) (samples []S, err error) {
	readableByteSize, err := getByteSizeToReachEOF(f)
	if err != nil {
		return nil, err
	}

	if readableByteSize == 0 {
		samples = make([]S, 0)
	} else if samplingN > 0 && readableByteSize > samplingN*samplingBytes {
		samples = make([]S, samplingN)
	} else {
		samples = make([]S, readableByteSize/samplingBytes)
	}
	err = binary.Read(f, binary.LittleEndian, samples)
	return
}

func getByteSizeToReachEOF(f io.ReadSeeker) (int64, error) {
	current, err := f.Seek(0, 1)
	if err != nil {
		return -1, err
	}
	end, err := f.Seek(0, 2)
	if err != nil {
		return -1, err
	}

	// Undo Seek Offset
	c, err := f.Seek(current, 0)
	if err != nil {
		return -1, err
	} else if current != c {
		return -1, errors.New("seek Offset could not be undone")
	}

	return end - current, nil
}
