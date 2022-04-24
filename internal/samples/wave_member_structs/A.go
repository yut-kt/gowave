// Code generated by gen/samples_generator.go; DO NOT EDIT.
package wave_member_structs

import (
	"fmt"
	"os"

	"github.com/yut-kt/gowave/internal/chunk"
)

func GetRiffChunkA() *chunk.RiffChunk {
	return &chunk.RiffChunk{
		ID:     "RIFF",
		Size:   1,
		Format: "WAVE",
	}
}

func GetFmtChunkA() *chunk.FmtChunk {
	return &chunk.FmtChunk{
		ID:            "fmt ",
		Size:          16,
		AudioFormat:   1,
		NumChannels:   1,
		SampleRate:    8000,
		ByteRate:      16000,
		BlockAlign:    2,
		BitsPerSample: 16,
	}
}

func GetDataChunkA() *chunk.DataChunk {
	f, err := os.Open("internal/samples/waves/A.wav")
	if err != nil {
		panic(err)
	}

	// Byte size required for SamplingData offset.
	// RiffChunkByte+(SubChunk1HeaderByte+SubChunk1Size)+SubChunk2HeaderByte
	byteSize := 12 + (8 + 16) + 8
	buf := make([]byte, byteSize)
	n, err := f.Read(buf)
	if n != byteSize {
		panic(fmt.Errorf("n(%v) != byteSize(%v)", n, byteSize))
	}
	if err != nil {
		panic(err)
	}

	return &chunk.DataChunk{
		File: f,
		ID:   "data",
		Size: 0,
		Data: nil,
	}
}