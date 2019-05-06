# Go module for appending files

[![Build Status](https://travis-ci.org/gonejack/gwriter.svg?branch=master)](https://travis-ci.org/gonejack/gwriter)
[![GoDoc](https://godoc.org/github.com/gonejack/gwriter?status.svg)](https://godoc.org/github.com/gonejack/gwriter)
[![Go Report Card](https://goreportcard.com/badge/github.com/gonejack/gwriter)](https://goreportcard.com/report/github.com/gonejack/gwriter)
[![GitHub license](https://img.shields.io/github/license/gonejack/gwriter.svg?color=blue)](LICENSE.md)

```go
func ExampleNewWriter() {
	conf := config.Config{
		PathTpl:  "{dir}/{filename}{base_ext}{write_ext}",
		BaseExt:  ".msg",
		WriteExt: "",
		PathInfo: map[string]string{
			"{dir}":      os.Getenv("DIR"),
			"{filename}": os.Getenv("FILENAME"),
		},
		UpdateMoment: "00:01:00",
	}

	writer := NewWriter("test-file", conf)
	writer.Start()
	writer.WriteString("this is string")
	writer.WriteBytes([]byte("this is bytes"))
	writer.Stop()
}
```