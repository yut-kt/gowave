package gowave

import (
	"io"

	"github.com/yut-kt/gowave/chunk"
)

type Wave struct {
	riffChunk *chunk.RiffChunk
	fmtChunk  *chunk.FmtChunk
	dataChunk *chunk.DataChunk
	file      io.Reader
}

func New(file io.Reader) (*Wave, error) {
	wave := &Wave{file: file}
	if err := wave.chunkRead(file); err != nil {
		return nil, err
	}
	return wave, nil
}

func (wave *Wave) chunkRead(file io.Reader) error {
	riffChunk, err := chunk.NewRiffChunk(file)
	if err != nil {
		return err
	}
	wave.riffChunk = riffChunk

	fmtChunk, err := chunk.NewFmtChunk(file)
	if err != nil {
		return err
	}
	wave.fmtChunk = fmtChunk

	dataChunk, err := chunk.NewDataChunk(file)
	if err != nil {
		return err
	}
	wave.dataChunk = dataChunk

	return nil
}

// ReadSamples is jhgkjhkjga
func (wave *Wave) ReadSamples() (interface{}, error) {
	const fullRead = -1
	data, err := wave.dataChunk.ReadData(wave.file, wave.fmtChunk.GetBitsPerSample(), fullRead)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (wave *Wave) ReadNSamples(samplingNum int) (interface{}, error) {
	data, err := wave.dataChunk.ReadData(wave.file, wave.fmtChunk.GetBitsPerSample(), samplingNum)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (wave *Wave) GetSamplesAlreadyRead() interface{} {
	return wave.dataChunk.GetData()
}

func (wave *Wave) GetNumChannels() uint16 {
	return wave.fmtChunk.GetNumChannels()
}

func (wave Wave) GetSampleRate() uint32 {
	return wave.fmtChunk.GetSampleRate()
}
