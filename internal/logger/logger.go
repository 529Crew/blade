package logger

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	TimeLocation *time.Location
)

func init() {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	TimeLocation = location
}

type writer struct {
	io.Writer
}

func (w writer) Write(b []byte) (n int, err error) {
	timestamp := time.Now().UTC().In(TimeLocation).Format("3:04:05.000Z PM | ") // Include milliseconds directly in the format
	return w.Writer.Write(append([]byte(timestamp), b...))
}

var (
	Log = log.New(&writer{os.Stdout}, "[BLADE] ", 0)
)
