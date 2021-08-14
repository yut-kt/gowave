package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

const SampleWaveDir = "../internal/samples/waves"
const SampleWaveFormatDir = "../internal/samples/format_jsons"
const SampleWaveInfoStructDir = "../internal/samples/wave_member_structs"

type waveInfo struct {
	ChunkID       string `json:"chunk_id"`
	ChunkSize     uint32 `json:"chunk_size"`
	Format        string `json:"format"`
	SubChunk1ID   string `json:"sub_chunk_1_id"`
	SubChunk1Size uint32 `json:"sub_chunk_1_size"`
	AudioFormat   uint16 `json:"audio_format"`
	NumChannels   uint16 `json:"num_channels"`
	SampleRate    uint32 `json:"sample_rate"`
	ByteRate      uint32 `json:"byte_rate"`
	BlockAlign    uint16 `json:"block_align"`
	BitsPerSample uint16 `json:"bits_per_sample"`
	SubChunk2ID   string `json:"sub_chunk_2_id"`
	SubChunk2Size uint32 `json:"sub_chunk_2_size"`
	FileName      string
}

//go:generate go run .
func main() {
	files, err := os.ReadDir(SampleWaveFormatDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// PathSetting
		readPath := filepath.Join(SampleWaveFormatDir, file.Name())
		nameExt := strings.Split(file.Name(), ".")
		fileName := strings.Join(nameExt[:len(nameExt)-1], ".")
		writeWavPath := filepath.Join(SampleWaveDir, fileName) + ".wav"
		writeGoPath := filepath.Join(SampleWaveInfoStructDir, fileName) + ".go"
		fmt.Println("Read:", readPath)

		// ReadFile
		b, err := os.ReadFile(readPath)
		if err != nil {
			panic(err)
		}
		var w waveInfo
		if err := json.Unmarshal(b, &w); err != nil {
			panic(err)
		}
		w.FileName = fileName

		// WriteSampleWave
		data := []interface{}{
			w.ChunkID, w.ChunkSize, w.Format,
			w.SubChunk1ID, w.SubChunk1Size, w.AudioFormat, w.NumChannels, w.SampleRate, w.ByteRate, w.BlockAlign, w.BitsPerSample,
			w.SubChunk2ID, w.SubChunk2Size,
		}
		buf := new(bytes.Buffer)
		for _, d := range data {
			switch v := d.(type) {
			case string:
				if err := binary.Write(buf, binary.BigEndian, []byte(v)); err != nil {
					panic(err)
				}
			case uint16:
				if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
					panic(err)
				}
			case uint32:
				if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
					panic(err)
				}
			default:
				panic("unknown type")
			}
		}
		if err := os.WriteFile(writeWavPath, buf.Bytes(), 0644); err != nil {
			panic(err)
		}
		fmt.Println("Write:", writeWavPath)

		// WriteSampleWaveInfoStruct
		buf = new(bytes.Buffer)
		if err := template.Must(template.New("prog").Parse(prog)).Execute(buf, w); err != nil {
			panic(err)
		}
		parsed, err := format.Source(buf.Bytes())
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(writeGoPath, parsed, 0644); err != nil {
			panic(err)
		}
	}
}

const prog = `
// Code generated by gen/samples_generator.go; DO NOT EDIT.
package wave_member_structs

import (
	"fmt"
	"io"
	"os"

	"github.com/yut-kt/gowave/internal/chunk"
)

func GetRiffChunk{{.FileName}}() *chunk.RiffChunk {
	return &chunk.RiffChunk{
		ID:     "{{.ChunkID}}",
		Size:   {{.ChunkSize}},
		Format: "{{.Format}}",
	}
}

func GetFmtChunk{{.FileName}}() *chunk.FmtChunk {
	return &chunk.FmtChunk{
		ID:            "{{.SubChunk1ID}}",
		Size:          {{.SubChunk1Size}},
		AudioFormat:   {{.AudioFormat}},
		NumChannels:   {{.NumChannels}},
		SampleRate:    {{.SampleRate}},
		ByteRate:      {{.ByteRate}},
		BlockAlign:    {{.BlockAlign}},
		BitsPerSample: {{.BitsPerSample}},
	}
}

func GetDataChunk{{.FileName}}() *chunk.DataChunk {
	return &chunk.DataChunk{
		ID:   "{{.SubChunk2ID}}",
		Size: {{.SubChunk2Size}},
		Data: nil,
	}
}

func GetFileA() io.ReadSeeker {
	f, err := os.Open("internal/samples/waves/{{.FileName}}.wav")
	if err != nil {
		panic(err)
	}

	// Byte size required for SamplingData offset.
	// RiffChunkByte+(SubChunk1HeaderByte+SubChunk1Size)+SubChunk2HeaderByte
	byteSize := 12 + (8 + {{.SubChunk1Size}}) + 8
	buf := make([]byte, byteSize)
	n, err := f.Read(buf)
	if n != byteSize {
		panic(fmt.Errorf("n(%v) != byteSize(%v)", n, byteSize))
	}
	if err != nil {
		panic(err)
	}

	return f
}

`
