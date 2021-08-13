package gowave

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/yut-kt/gowave/internal/samples/wave_member_structs"

	"github.com/yut-kt/gowave/internal/chunk"
)

const aWaveFilePath = "internal/samples/waves/A.wav"

var (
	waveA = &Wave{
		riffChunk: wave_member_structs.GetRiffChunkA(),
		fmtChunk:  wave_member_structs.GetFmtChunkA(),
		dataChunk: wave_member_structs.GetDataChunkA(),
		file:      wave_member_structs.GetFileA(),
	}
)

// Equal test of io is difficult, so check only chunk.
func notWaveEqual(gotWave *Wave, wantWave *Wave) bool {
	return !reflect.DeepEqual(gotWave.riffChunk, wantWave.riffChunk) ||
		!reflect.DeepEqual(gotWave.fmtChunk, wantWave.fmtChunk) ||
		!reflect.DeepEqual(gotWave.dataChunk, wantWave.dataChunk)
}

func TestNew(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Wave
		wantErr bool
	}{
		{
			name:    "A",
			args:    args{filePath: aWaveFilePath},
			want:    waveA,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.filePath)
			if err != nil {
				t.Error(err)
			}
			got, err := New(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if notWaveEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWave_GetNumChannels(t *testing.T) {
	type fields struct {
		fmtChunk *chunk.FmtChunk
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name:   "A",
			fields: fields{fmtChunk: wave_member_structs.GetFmtChunkA()},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wave := &Wave{
				fmtChunk: tt.fields.fmtChunk,
			}
			if got := wave.GetNumChannels(); got != tt.want {
				t.Errorf("GetNumChannels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWave_GetSampleRate(t *testing.T) {
	type fields struct {
		fmtChunk *chunk.FmtChunk
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "A",
			fields: fields{fmtChunk: wave_member_structs.GetFmtChunkA()},
			want:   8000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wave := &Wave{
				fmtChunk: tt.fields.fmtChunk,
			}
			if got := wave.GetSampleRate(); got != tt.want {
				t.Errorf("GetSampleRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWave_GetSamplesAlreadyRead(t *testing.T) {
	type fields struct {
		riffChunk *chunk.RiffChunk
		fmtChunk  *chunk.FmtChunk
		dataChunk *chunk.DataChunk
		file      io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wave := &Wave{
				riffChunk: tt.fields.riffChunk,
				fmtChunk:  tt.fields.fmtChunk,
				dataChunk: tt.fields.dataChunk,
				file:      tt.fields.file,
			}
			if got := wave.GetSamplesAlreadyRead(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSamplesAlreadyRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWave_ReadNSamples(t *testing.T) {
	type args struct {
		samplingNum int
	}
	tests := []struct {
		name    string
		args    args
		wave    *Wave
		want    interface{}
		wantErr bool
	}{
		{
			name:    "num error",
			args:    args{0},
			wave:    waveA,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.wave.ReadNSamples(tt.args.samplingNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNSamples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadNSamples() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWave_ReadSamples(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name     string
		args     args
		wantWave *Wave
		wantErr  bool
	}{
		{name: "testA", args: args{filePath: aWaveFilePath}, wantWave: waveA, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.filePath)
			if err != nil {
				t.Errorf("os.Open(%v) error = %v", tt.args.filePath, err)
			}

			gotWave, err := New(f)
			if err != nil {
				t.Errorf("New(f) error = %v", err)
			}

			got, err := gotWave.ReadSamples()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadSamples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want, err := tt.wantWave.ReadSamples()
			if err != nil {
				t.Errorf("tt.wantWave.ReadSamples() error = %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("ReadSamples() got = %v, want %v", got, want)
			}
		})
	}
}

func TestWave_chunkRead(t *testing.T) {
	type fields struct {
		riffChunk *chunk.RiffChunk
		fmtChunk  *chunk.FmtChunk
		dataChunk *chunk.DataChunk
		file      io.Reader
	}
	type args struct {
		file io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wave := &Wave{
				riffChunk: tt.fields.riffChunk,
				fmtChunk:  tt.fields.fmtChunk,
				dataChunk: tt.fields.dataChunk,
				file:      tt.fields.file,
			}
			if err := wave.chunkRead(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("chunkRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
