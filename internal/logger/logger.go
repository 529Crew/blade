package logger

import (
	"io"
	"log"
	"os"
	"time"
)

type writer struct {
	io.Writer
	timeFormat string
}

func (w writer) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(w.timeFormat)), b...))
}

var (
	Log = log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, "[BLADE] ", 0)
)
