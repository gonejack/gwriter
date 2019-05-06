package basicWriter

import (
	"bufio"
	"fmt"
	"github.com/gonejack/glogger"
	"github.com/gonejack/gwriter/config"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type writer struct {
	logger glogger.Logger
	config config.Config

	stringInput chan string
	bytesInput  chan []byte

	signal   chan os.Signal
	fp       *os.File
	writer   *bufio.Writer
	wroteLen int64

	flushTimer <-chan time.Time
	closeTimer <-chan time.Time
}

func (w *writer) mainRoutine() {
	stopping := false

	for {
		select {
		case <-w.signal:
			stopping = true
		default:
			if stopping && len(w.stringInput) == 0 && len(w.bytesInput) == 0 {
				w.close()

				return
			} else {
				select {
				case <-w.flushTimer:
					w.flush()
				case <-w.closeTimer:
					w.close()
				case s := <-w.stringInput:
					w.write(s, nil)
				case b := <-w.bytesInput:
					w.write("", b)
				}
			}
		}
	}
}

func (w *writer) WriteString(s string) {
	w.stringInput <- s
}
func (w *writer) WriteBytes(bs []byte) {
	w.bytesInput <- bs
}
func (w *writer) Start() {
	w.logger.Infof("开始启动")

	go w.mainRoutine()

	w.logger.Infof("启动完成")
}
func (w *writer) Stop() {
	w.logger.Infof("开始关闭")

	w.signal <- os.Interrupt
	for w.writer != nil {
		time.Sleep(time.Millisecond * 100)
	}

	w.logger.Infof("关闭完成")
}

func (w *writer) open() bool {
	if w.fp == nil {
		if p := w.genPath(); p != "" {
			fp, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err == nil {
				w.logger.Infof("创建文件[%s]", p)

				w.fp = fp
				w.writer = bufio.NewWriter(fp)
				w.flushTimer = time.Tick(time.Second)
				w.closeTimer = time.After(time.Duration(w.calSleepSec()) * time.Second)

				if w.config.UpdateSize > 0 {
					w.logger.Infof("文件计划关闭大小为: %s bytes", w.config.UpdateSize)
				}
			} else {
				w.logger.Errorf("创建文件[%s]失败: %s", p, err)
			}
		}
	}

	return w.writer != nil
}
func (w *writer) flush() {
	if w.writer != nil {
		err := w.writer.Flush()

		if err != nil {
			w.logger.Errorf("写内容失败：%s", err)
		}
	}
}
func (w *writer) close() {
	if w.fp != nil {
		w.flush()

		err := w.fp.Close()
		if err != nil {
			w.logger.Errorf("关闭文件出错: %s", err)
		}

		w.logger.Infof("关闭文件[%s]", w.fp.Name())

		if w.config.WriteExt != "" {
			w.removeWriteExt()
		}
	}

	w.fp = nil
	w.writer = nil
	w.flushTimer = nil
	w.closeTimer = nil
	w.wroteLen = 0
}
func (w *writer) write(str string, bs []byte) {
	if w.open() {
		var length int

		if length = len(str); length > 0 {
			_, err := w.writer.WriteString(str)
			if err != nil {
				w.logger.Errorf("写内容失败: %s", err)
			}

			w.wroteLen += int64(length)
		}

		if length = len(bs); length > 0 {
			_, err := w.writer.Write(bs)
			if err != nil {
				w.logger.Errorf("写内容失败: %s", err)
			}

			w.wroteLen += int64(length)
		}

		if w.config.UpdateSize > 0 && w.wroteLen >= w.config.UpdateSize {
			w.logger.Infof("写入长度达到限制[%s bytes]", w.config.UpdateSize)

			w.close()
		}
	} else {
		w.logger.Errorf("没有打开的文件，写入失败")
	}
}
func (w *writer) genPath() (path string) {
	now := time.Now()

	info := map[string]string{
		"{date}":      now.Format("20060102"),
		"{ts}":        strconv.FormatInt(now.Unix(), 10),
		"{base_ext}":  w.config.BaseExt,
		"{write_ext}": w.config.WriteExt,
	}

	for macro, replacement := range w.config.PathInfo {
		info[macro] = replacement
	}

	var replaces []string
	for macro, replacement := range info {
		replaces = append(replaces, macro, replacement)
	}

	path = strings.NewReplacer(replaces...).Replace(w.config.PathTpl)
	dir := filepath.Dir(path)

	if _, e := os.Stat(dir); os.IsNotExist(e) {
		if e = os.MkdirAll(dir, 0755); e == nil {
			w.logger.Infof("创建文件夹[%s]", dir)
		} else {
			w.logger.Errorf("创建文件夹[%s]失败: %s", dir, e)

			path = ""
		}
	}

	return
}
func (w *writer) removeWriteExt() {
	if o := w.fp.Name(); strings.HasSuffix(o, w.config.WriteExt) {
		n := strings.TrimSuffix(o, w.config.WriteExt)
		e := os.Rename(o, n)

		if e == nil {
			w.logger.Infof("重命名文件[%s => %s]", o, n)
		} else {
			w.logger.Infof("重命名文件[%s => %s]出错: %s", o, n, e)
		}
	}
}
func (w *writer) calSleepSec() (sec int) {
	now := time.Now()

	start := now
	if w.config.UpdateMoment != "" {
		todayMoment := fmt.Sprintf("%s %s", now.Format("2006-01-02"), w.config.UpdateMoment)

		if moment, e := time.ParseInLocation("2006-01-02 15:04:05", todayMoment, now.Location()); e == nil {
			start = moment
		} else {
			w.logger.Errorf("配置的文件更新时刻[%s]解析出错: %s，缺省为现在时刻", w.config.UpdateMoment, e)
		}
	}

	next := start
	for next.Before(now) || next.Equal(now) {
		if w.config.UpdatePeriod > 0 {
			next = next.Add(time.Duration(w.config.UpdatePeriod) * time.Second)
		} else {
			next = next.Add(time.Hour * 24)
		}
	}

	sec = int(math.Ceil(next.Sub(now).Seconds()))

	w.logger.Infof("文件计划关闭时间为: %s", next.Format("2006-01-02 15:04:05"))

	return
}

// New create an instance of basic Writer
func New(name string, config config.Config) (w *writer) {
	w = &writer{
		logger: glogger.NewLogger(name),

		config: config,

		stringInput: make(chan string, 100),
		bytesInput:  make(chan []byte, 100),
		signal:      make(chan os.Signal),
	}

	return
}
