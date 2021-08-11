# gowave

[![v0.1.0](https://img.shields.io/github/v/release/yut-kt/gowave?logoColor=ff69b4&style=social)]()
[![coverage](https://img.shields.io/badge/coverage-72.4%25-green)](https://raw.githubusercontent.com/yut-kt/gowave/main/coverage/v0.1.0)  
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

## Contribution
1. Fork ([https://github.com/yut-kt/gowave/fork](https://github.com/yut-kt/gowave/fork))
2. Checkout the latest version of branch
3. Create a feature branch
4. Commit your changes
5. Run test suite with the `make test` command and confirm that it passes
6. Create new Pull Request

## License
gowave is released under the [MIT License](https://raw.githubusercontent.com/yut-kt/gowave/main/LICENSE).
