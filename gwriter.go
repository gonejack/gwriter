/*
Package writer is an file library for appending file.
*/
package gwriter

import (
	"github.com/gonejack/gwriter/config"
	"github.com/gonejack/gwriter/internal/basicWriter"
	"os"
)

// NewWriter creates a new writer instance of basicWriter match Writer interface
func NewWriter(name string, config config.Config) (w Writer) {
	return basicWriter.New(name, config)
}

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

	writer := NewWriter("心跳消息文件", conf)
	writer.Start()
	writer.Stop()
}
