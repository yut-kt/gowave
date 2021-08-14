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
func (chunk *DataChunk) ReadData(bitsPerSample uint16, samplingNum int64) (interface{}, error) {
	data, err := chunk.readSamples(bitsPerSample, samplingNum)
	if err != nil {
		return nil, err
	}

	switch v := chunk.Data.(type) {
	case []uint8:
		d, ok := data.([]uint8)
		if !ok {
			return nil, errors.New("not match interface type")
		}
		chunk.Data = append(v, d...)
	case []int16:
		d, ok := data.([]int16)
		if !ok {
			return nil, errors.New("not match interface type")
		}
		chunk.Data = append(v, d...)
	default:
		chunk.Data = data
	}
	return data, nil
}

func (chunk *DataChunk) readSamples(bitsPerSample uint16, samplingN int64) (interface{}, error) {
	var data interface{}

	readableByteSize, err := chunk.getByteSizeToReachEOF()
	if err != nil {
		return nil, err
	}

	switch bitsPerSample {
	case 8:
		const samplingBytes = 1
		if readableByteSize == 0 {
			return make([]uint8, 0), nil
		}
		if samplingN > 0 && readableByteSize > samplingN*samplingBytes {
			data = make([]uint8, samplingN)
		} else {
			data = make([]uint8, readableByteSize/samplingBytes)
		}
	case 16:
		const samplingBytes = 2
		if readableByteSize == 0 {
			return make([]int16, 0), nil
		}
		if samplingN > 0 && readableByteSize > samplingN*samplingBytes {
			data = make([]int16, samplingN)
		} else {
			data = make([]int16, readableByteSize/samplingBytes)
		}
	default:
		return nil, errors.New("not supported bitPerSample number")
	}

	if err := binary.Read(chunk.File, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (chunk *DataChunk) getByteSizeToReachEOF() (int64, error) {
	current, err := chunk.File.Seek(0, 1)
	if err != nil {
		return -1, err
	}
	end, err := chunk.File.Seek(0, 2)
	if err != nil {
		return -1, err
	}

	// Undo Seek Offset
	c, err := chunk.File.Seek(current, 0)
	if err != nil {
		return -1, err
	} else if current != c {
		return -1, errors.New("seek Offset could not be undone")
	}

	return end - current, nil
}
