package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const WaveSampleDir = "../internal/samples/waves"
const WaveSampleFormatDir = "../internal/samples/format_jsons"

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
}

//go:generate go run .
func main() {
	files, err := os.ReadDir(WaveSampleFormatDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		readPath := filepath.Join(WaveSampleFormatDir, file.Name())
		nameExt := strings.Split(file.Name(), ".")
		writePath := filepath.Join(WaveSampleDir, strings.Join(nameExt[:len(nameExt)-1], ".")) + ".wav"
		fmt.Println("Read:", readPath)

		b, err := os.ReadFile(readPath)
		if err != nil {
			panic(err)
		}

		var w waveInfo
		if err := json.Unmarshal(b, &w); err != nil {
			panic(err)
		}
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
		if err := os.WriteFile(writePath, buf.Bytes(), os.FileMode(0644)); err != nil {
			panic(err)
		}
		fmt.Println("Write:", writePath)
	}
}
