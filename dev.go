package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

var colors = [...]int{
	Debug: gray,
	Info:  blue,
	Warn:  yellow,
	Error: red,
}

var strings = [...]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
}

type devHandler struct {
	writer io.Writer
	start  time.Time
	sync.Mutex
}

func NewDevHandler(w io.Writer) Handler {
	return &devHandler{
		writer: w,
		start:  clock(),
	}
}

func (d *devHandler) Handle(e *Entry) {
	color := colors[e.Level]
	level := strings[e.Level]

	d.Lock()
	defer d.Unlock()

	elapsed := time.Since(d.start) / time.Second
	fmt.Fprintf(d.writer, "\033[%dm%6s\033[0m[%04d] %-25s", color, level, elapsed, e.Message)

	for k, v := range e.Fields {
		fmt.Fprintf(d.writer, " \033[%dm%s\033[0m=%v", color, k, v)
	}

	fmt.Fprintln(d.writer)
}
