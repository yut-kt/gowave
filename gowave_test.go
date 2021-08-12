package gowave

import (
	"io"
	"os"
	"reflect"
	"testing"

	sample_wave_structs "github.com/yut-kt/gowave/internal/samples/wave_member_structs"

	"github.com/yut-kt/gowave/internal/chunk"
)

const aWaveFilePath = "internal/samples/waves/A.wav"

func TestNew(t *testing.T) {
	type args struct {
		file io.Reader
	}
	aFile, err := os.Open(aWaveFilePath)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		args    args
		want    *Wave
		wantErr bool
	}{
		{
			name: "A",
			args: args{file: aFile},
			want: &Wave{
				riffChunk: sample_wave_structs.GetRiffChunkA(),
				fmtChunk:  sample_wave_structs.GetFmtChunkA(),
				dataChunk: sample_wave_structs.GetDataChunkA(),
				file:      aFile,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
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
			fields: fields{fmtChunk: sample_wave_structs.GetFmtChunkA()},
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
			fields: fields{fmtChunk: sample_wave_structs.GetFmtChunkA()},
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
	type fields struct {
		riffChunk *chunk.RiffChunk
		fmtChunk  *chunk.FmtChunk
		dataChunk *chunk.DataChunk
		file      io.Reader
	}
	type args struct {
		samplingNum int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
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
			got, err := wave.ReadNSamples(tt.args.samplingNum)
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
	type fields struct {
		riffChunk *chunk.RiffChunk
		fmtChunk  *chunk.FmtChunk
		dataChunk *chunk.DataChunk
		file      io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
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
			got, err := wave.ReadSamples()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadSamples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadSamples() got = %v, want %v", got, tt.want)
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
