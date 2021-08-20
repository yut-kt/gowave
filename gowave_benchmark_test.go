package gowave

import (
	"os"
	"testing"

	"github.com/yut-kt/gowave/internal/samples/wave_member_structs"
)

var waveFiles = map[string]*Wave{
	"internal/samples/waves/X.wav": {
		riffChunk: wave_member_structs.GetRiffChunkX(),
		fmtChunk:  wave_member_structs.GetFmtChunkX(),
		dataChunk: wave_member_structs.GetDataChunkX(),
	},
}

func BenchmarkNew(b *testing.B) {
	for waveFile := range waveFiles {
		f, err := os.Open(waveFile)
		if err != nil {
			b.Errorf("os.Open(%v) error = %v", waveFile, err)
		}
		if _, err := New(f); err != nil {
			b.Errorf("New(f) error = %v", err)
		}
	}
}

func BenchmarkWave_ReadSamples(b *testing.B) {
	for _, wave := range waveFiles {
		if _, err := wave.ReadSamples(); err != nil {
			b.Errorf("wave.ReadSamples() error = %v", err)
		}
	}
}

func BenchmarkWave_ReadNSamples(b *testing.B) {
	for _, wave := range waveFiles {
		if _, err := wave.ReadNSamples(1000); err != nil {
			b.Errorf("wave.ReadNSamples(1000) error = %v", err)
		}
	}
}
