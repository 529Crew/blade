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

func newLogWriter() io.Writer {
	/* create logs folder if it doesnt exist */
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	/* create log file for current date */
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logFilePath := logDir + "/" + logFileName
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return &writer{multiWriter}
}

var (
	Log = log.New(newLogWriter(), "[BLADE] ", 0)
)
