// Package gowave provides support for reading WAV files.
package gowave

import (
	"errors"
	"io"

	"github.com/yut-kt/gowave/internal/chunk"
)

// Wave is a structure that handles wav files.
type Wave struct {
	riffChunk *chunk.RiffChunk
	fmtChunk  *chunk.FmtChunk
	dataChunk *chunk.DataChunk
}

// New is a function to construct Wave struct.
func New(file io.ReadSeeker) (*Wave, error) {
	wave := new(Wave)
	if err := wave.chunkRead(file); err != nil {
		return nil, err
	}
	return wave, nil
}

func (wave *Wave) chunkRead(file io.ReadSeeker) error {
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

// ReadSamples is a function to read all samples wave data.
func (wave *Wave) ReadSamples() (interface{}, error) {
	const fullRead = -1
	data, err := wave.dataChunk.ReadData(wave.fmtChunk.BitsPerSample, fullRead)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadNSamples is a function to read N samples wave data.
func (wave *Wave) ReadNSamples(samplingNum int64) (interface{}, error) {
	if samplingNum < 1 {
		return nil, errors.New("samplingNum is only natural number")
	}
	data, err := wave.dataChunk.ReadData(wave.fmtChunk.BitsPerSample, samplingNum)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetSamplesAlreadyRead is a function to get already read samples wave data.
func (wave *Wave) GetSamplesAlreadyRead() interface{} {
	return wave.dataChunk.Data
}

// GetNumChannels is a function to get num channels.
func (wave *Wave) GetNumChannels() uint16 {
	return wave.fmtChunk.NumChannels
}

// GetSampleRate is a function to get sample rate.
func (wave *Wave) GetSampleRate() uint32 {
	return wave.fmtChunk.SampleRate
}
