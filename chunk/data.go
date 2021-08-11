package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"
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

	chunk := &DataChunk{
		id:   string(chunkHeaderBytes[:4]),
		size: binary.LittleEndian.Uint32(chunkHeaderBytes[4:]),
	}
	if err := chunk.validate(); err != nil {
		return nil, err
	}
	return chunk, nil
}

func (chunk *DataChunk) validate() error {
	if chunk.id != "data" {
		return errors.New("DataChunk: ChunkID must be [data]")
	}
	return nil
}

func (chunk *DataChunk) ReadData(file io.Reader, bitsPerSample uint16, samplingNum int) (interface{}, error) {
	var (
		data    interface{}
		funcErr error
	)
	if samplingNum == -1 {
		data, funcErr = readAllSamples(file, bitsPerSample)
		if funcErr != nil {
			return nil, funcErr
		}
	} else {
		data, funcErr = readNSamples(file, bitsPerSample, samplingNum)
		if funcErr != nil {
			return nil, funcErr
		}
	}

	switch v := chunk.data.(type) {
	case []uint8:
		d, ok := data.([]uint8)
		if !ok {
			return nil, errors.New("not match interface type")
		}
		chunk.data = append(v, d...)
	case []int16:
		d, ok := data.([]int16)
		if !ok {
			return nil, errors.New("not match interface type")
		}
		chunk.data = append(v, d...)
	default:
		chunk.data = data
	}
	return data, nil
}

func readAllSamples(file io.Reader, bitsPerSample uint16) (interface{}, error) {
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

func readNSamples(file io.Reader, bitsPerSample uint16, samplingN int) (interface{}, error) {
	var (
		data               interface{}
		buf                *bytes.Reader
		wasReadableSampleN int
	)

	switch bitsPerSample {
	case 8:
		const samplingBytes = 1
		b := make([]byte, samplingN*samplingBytes)
		n, err := io.ReadFull(file, b)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return nil, err
		}
		// store
		buf = bytes.NewReader(b[:n])
		wasReadableSampleN = n / samplingBytes
		data = make([]uint8, wasReadableSampleN)
	case 16:
		const samplingBytes = 2
		b := make([]byte, samplingN*samplingBytes)
		n, err := io.ReadFull(file, b)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return nil, err
		}
		// store
		buf = bytes.NewReader(b[:n])
		wasReadableSampleN = n / samplingBytes
		data = make([]int16, wasReadableSampleN)
	default:
		return nil, errors.New("not supported bitPerSample number")
	}
	if err := binary.Read(buf, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (chunk *DataChunk) GetData() interface{} {
	return chunk.data
}
