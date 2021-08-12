# gowave

[![v0.1.0](https://img.shields.io/github/v/release/yut-kt/gowave?logoColor=ff69b4&style=social)]()
[![Test](https://github.com/yut-kt/gowave/actions/workflows/default_branch_test.yaml/badge.svg)](https://github.com/yut-kt/gowave/actions/workflows/default_branch_test.yaml)
[![coverage](https://img.shields.io/badge/coverage-72.4%25-green)](https://raw.githubusercontent.com/yut-kt/gowave/main/coverage/v0.1.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/yut-kt/gowave)](https://goreportcard.com/report/github.com/yut-kt/gowave)  
[![Go Reference](https://pkg.go.dev/badge/github.com/yut-kt/gowave.svg)](https://pkg.go.dev/github.com/yut-kt/gowave)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/yut-kt/gowave/main/LICENSE)


**Wave file read support for Go language**

## Install
```bash
$ go get github.com/yut-kt/gowave
```

## Import
```go
import (
    "github.com/yut-kt/gohowave"
)
```

## Usage
```go
import (
    "fmt"
    "os"
    
    "github.com/yut-kt/gowave"
)

func main() {
    const WaveFile = "XXX.wav"
	
    // Open File
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
        fmt.Println(len(v)) // 100000lengthSamples
    }
}
```

See [gowave_examples_test.go](https://github.com/yut-kt/gowave/blob/main/gowave_examples_test.go) for detailed Usage

## Supported format
Format

- PCM
- ~~IEEE float (read-only)~~

Number of channels

- 1(mono)
- ~~2(stereo)~~

Bits per sample

- ~~32-bit~~
- ~~24-bit~~
- 16-bit
- 8-bit

## License
gowave is released under the [MIT License](https://raw.githubusercontent.com/yut-kt/gowave/main/LICENSE).
