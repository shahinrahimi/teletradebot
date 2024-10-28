package logger

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type Logger struct {
	Log     *log.Logger
	logFile *os.File
}

func New(name string, saveToDisk bool) *Logger {
	timeNow := time.Now().Format("2006-01-02")
	lowerName := strings.ToLower(name)
	upperName := strings.ToUpper(name)
	if saveToDisk {
		if err := os.MkdirAll("logs", 0755); err != nil {
			log.Fatal(err)
		}
		filename := path.Join("logs", lowerName+"-"+timeNow+".txt")
		logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		wr := io.MultiWriter(os.Stdout, logFile)
		return &Logger{
			Log:     log.New(wr, "["+upperName+"] ", log.LstdFlags),
			logFile: logFile,
		}
	} else {
		return &Logger{
			Log: log.New(os.Stdout, "["+upperName+"] ", log.LstdFlags),
		}
	}

}

func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
