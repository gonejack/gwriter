package gwriter

import (
	"github.com/gonejack/gwriter/config"
	"os"
	"testing"
	"time"
)

func TestNewWriter(t *testing.T) {
	conf := config.Config{
		PathTpl:  "{dir}/{filename}{base_ext}{write_ext}",
		BaseExt:  ".msg",
		WriteExt: "",
		PathInfo: map[string]string{
			"{dir}":      ".",
			"{filename}": "testFile",
		},
		UpdateMoment: "00:01:00",
	}

	writer := NewWriter("writerTest", conf)

	writer.Start()
	writer.WriteString("this is string")
	writer.WriteBytes([]byte("this is bytes"))
	writer.Stop()

	time.Sleep(time.Second)
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
	writer.WriteString("this is string")
	writer.WriteBytes([]byte("this is bytes"))
	writer.Stop()
}
