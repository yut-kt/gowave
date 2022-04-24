package gowave

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/yut-kt/gowave/internal/samples/wave_member_structs"

	"github.com/yut-kt/gowave/internal/chunk"
)

type testWave int

const (
	testWaveA testWave = iota
)

// Equal test of io is difficult, so check only field excluding io.
func notWaveEqual(gotWave *Wave, wantWave *Wave) bool {
	return !reflect.DeepEqual(gotWave.riffChunk, wantWave.riffChunk) ||
		!reflect.DeepEqual(gotWave.fmtChunk, wantWave.fmtChunk) ||
		!reflect.DeepEqual(
			[]interface{}{gotWave.dataChunk.ID, gotWave.dataChunk.Size, gotWave.dataChunk.Data},
			[]interface{}{wantWave.dataChunk.ID, wantWave.dataChunk.Size, wantWave.dataChunk.Data})
}

func newTestWave(wave testWave) *Wave {
	switch wave {
	case testWaveA:
		return &Wave{
			riffChunk: wave_member_structs.GetRiffChunkA(),
			fmtChunk:  wave_member_structs.GetFmtChunkA(),
			dataChunk: wave_member_structs.GetDataChunkA(),
		}
	}
	return nil
}

func testPath(wave testWave) string {
	switch wave {
	case testWaveA:
		return "internal/samples/waves/A.wav"
	}
	return ""
}

// Wave initialization common function for testing
func newReadWaveFile(t *testing.T, filePath string) *Wave {
	f, err := os.Open(filePath)
	if err != nil {
		t.Errorf("os.Open(%v) error = %v", filePath, err)
	}

	w, err := New(f)
	if err != nil {
		t.Errorf("New(f) error = %v", err)
	}

	return w
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
			args:    args{filePath: testPath(testWaveA)},
			want:    newTestWave(testWaveA),
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
				fmt.Println(got.dataChunk.Data, tt.want.dataChunk.Data)
				fmt.Printf("%T\n", got.dataChunk.Data)
				fmt.Printf("%T\n", tt.want.dataChunk.Data)

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
	type args struct {
		filePath string
		readN    int64
	}
	tests := []struct {
		name     string
		args     args
		wantWave *Wave
		want     interface{}
	}{
		{name: "testA0", args: args{filePath: testPath(testWaveA), readN: 0}, wantWave: newTestWave(testWaveA)},
		{name: "testA5", args: args{filePath: testPath(testWaveA), readN: 5}, wantWave: newTestWave(testWaveA)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWave := newReadWaveFile(t, tt.args.filePath)

			if tt.args.readN > 0 {
				if _, err := gotWave.ReadNSamples(tt.args.readN); err != nil {
					t.Errorf("gotWave.ReadNSamples(%v) error = %v", tt.args.readN, err)
				}
				if _, err := tt.wantWave.ReadNSamples(tt.args.readN); err != nil {
					t.Errorf("tt.wantWave.ReadNSamples(%v) error = %v", tt.args.readN, err)
				}
			}

			if got := gotWave.GetSamplesAlreadyRead(); !reflect.DeepEqual(got, tt.wantWave.GetSamplesAlreadyRead()) {
				t.Errorf("GetSamplesAlreadyRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWave_ReadNSamples(t *testing.T) {
	type args struct {
		samplingNum int64
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
			wave:    newTestWave(testWaveA),
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
		{name: "testA", args: args{filePath: testPath(testWaveA)}, wantWave: newTestWave(testWaveA), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWave := newReadWaveFile(t, tt.args.filePath)

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
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "testA", args: args{filePath: testPath(testWaveA)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.filePath)
			if err != nil {
				t.Errorf("os.Open(%v) error = %v", tt.args.filePath, err)
			}

			wave := new(Wave)
			if err := wave.chunkRead(file); (err != nil) != tt.wantErr {
				t.Errorf("chunkRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
