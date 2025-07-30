package logger

import (
	"io"
	"log"
	"os"
)

var (
	logFile *os.File
)

func InitLoggerFile(appLogPath string) error {
	var err error
	logFile, err = os.OpenFile(appLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Redirect standard logger to both stdout and file
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return nil
}

func CloseLogFile() {
	if logFile != nil {
		_ = logFile.Close()
	}
}
