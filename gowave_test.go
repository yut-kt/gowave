package gowave_test

import (
	"fmt"
	"os"
	"reflect"

	"github.com/yut-kt/gowave"
)

const WaveFile = "internal/samples/waves/jvs001_VOICEACTRESS100_001.wav"

func Example() {
	var a, b, c, d, e, x, y int

	file, err := os.Open(WaveFile)
	if err != nil {
		panic(err)
	}

	// Initialization
	wave, err := gowave.New(file)
	if err != nil {
		panic(err)
	}

	// Read 100000 Samples
	samples, err := wave.ReadNSamples(100000)
	if err != nil {
		panic(err)
	}
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		a = len(v)
		fmt.Println("a:", a)
	}

	// Read 1000 Samples
	samples, err = wave.ReadNSamples(1000)
	if err != nil {
		panic(err)
	}
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		b = len(v)
		fmt.Println("b:", b)
	}

	// Returns the stock of Samples read so far
	samples = wave.GetSamplesAlreadyRead()
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		x = len(v)
		fmt.Println("x:", x)
	}

	fmt.Println("a+b == x:", a+b == x)

	// If the number of readable samples is exceeded,
	// the Samples up to EOF are returned instead of the specified Samples
	samples, err = wave.ReadNSamples(200000)
	if err != nil {
		panic(err)
	}
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		c = len(v)
		fmt.Println("c:", c)
	}

	// Returns [] if there is no Readable sample
	samples, err = wave.ReadSamples()
	if err != nil {
		panic(err)
	}
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		d = len(v)
		fmt.Println("d:", d)
	}

	// Returns [] if there is no Readable sample
	samples, err = wave.ReadSamples()
	if err != nil {
		panic(err)
	}
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		e = len(v)
		fmt.Println("e:", e)
	}

	samples = wave.GetSamplesAlreadyRead()
	switch v := samples.(type) {
	case []uint8:
	case []int16:
		y = len(v)
		fmt.Println("y:", y)
	}

	fmt.Println("a+b+c+d+e == y:", a+b+c+d+e == y)
	// Output:
	// a: 100000
	// b: 1000
	// x: 101000
	// a+b == x: true
	// c: 105905
	// d: 0
	// e: 0
	// y: 206905
	// a+b+c+d+e == y: true
}

func ExampleWave_GetSamplesAlreadyRead() {
	file, err := os.Open(WaveFile)
	if err != nil {
		panic(err)
	}

	// Initialization
	wave, err := gowave.New(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("a:", wave.GetSamplesAlreadyRead())

	samples, err := wave.ReadNSamples(5)
	if err != nil {
		panic(err)
	}
	fmt.Println("b:", reflect.DeepEqual(samples, wave.GetSamplesAlreadyRead()))

	_, err = wave.ReadNSamples(5)
	if err != nil {
		panic(err)
	}
	fmt.Println("c:", 10 == reflect.ValueOf(wave.GetSamplesAlreadyRead()).Len())

	// Output:
	// a: <nil>
	// b: true
	// c: true
}
