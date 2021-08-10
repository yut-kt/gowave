package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type DataChunk struct {
	id   string
	size uint32
	data interface{}
}

func NewDataChunk(file io.Reader) (*DataChunk, error) {
	const chunkHeaderByteSize = 8
	chunkHeaderBytes := make([]byte, chunkHeaderByteSize)
	if _, err := io.ReadFull(file, chunkHeaderBytes); err != nil {
		return nil, err
	}

	return &DataChunk{
		id:   string(chunkHeaderBytes[:4]),
		size: binary.LittleEndian.Uint32(chunkHeaderBytes[4:]),
	}, nil
}

func (dataChunk *DataChunk) ReadData(file io.Reader, bitsPerSample uint16, samplingNum int) (interface{}, error) {
	readAllFunc := func(file io.Reader, bitsPerSample uint16) (interface{}, error) {
		var data interface{}

		// TODO: Consider calculating from the seek position and the eof position in io.ReadSeeker.
		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewReader(b)

		switch bitsPerSample {
		case 8:
			const samplingBytes = 1
			if len(b) == 0 {
				return make([]uint8, 0), nil
			}
			data = make([]uint8, len(b)/samplingBytes)
		case 16:
			const samplingBytes = 2
			if len(b) == 0 {
				return make([]int16, 0), nil
			}
			data = make([]int16, len(b)/samplingBytes)
		default:
			return nil, errors.New("not supported bitPerSample number")
		}

		if err := binary.Read(buf, binary.LittleEndian, data); err != nil {
			return nil, err
		}
		return data, nil
	}

	readNumFunc := func(file io.Reader, bitsPerSample uint16, samplingNum int) (interface{}, error) {
		var data interface{}

		switch bitsPerSample {
		case 8:
			data = make([]uint8, samplingNum)
		case 16:
			data = make([]int16, samplingNum)
		default:
			return nil, errors.New("not supported bitPerSample number")
		}
		if err := binary.Read(file, binary.LittleEndian, data); err != nil {
			return nil, err
		}
		return data, nil
	}

	var data interface{}
	var funcErr error
	if samplingNum == -1 {
		data, funcErr = readAllFunc(file, bitsPerSample)
		if funcErr != nil {
			return nil, funcErr
		}
	} else {
		data, funcErr = readNumFunc(file, bitsPerSample, samplingNum)
		if funcErr != nil {
			return nil, funcErr
		}
	}

	switch v := dataChunk.data.(type) {
	case []uint8:
		fmt.Println("uint8")
		d, ok := data.([]uint8)
		if !ok {
			return nil, errors.New("not match interface type")
		}
		dataChunk.data = append(v, d...)
	case []int16:
		d, ok := data.([]int16)
		if !ok {
			return nil, errors.New("not match interface type")
		}
		dataChunk.data = append(v, d...)
	default:
		dataChunk.data = data
	}
	return data, nil
}

func (dataChunk *DataChunk) GetData() interface{} {
	return dataChunk.data
}
