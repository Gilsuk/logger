package logger

import (
	"os"
	"sync"

	"github.com/kpango/glg"
)

type Logger interface {
	// Close is safe to call multiple times
	Close() error
	Debug(format string, v ...interface{})
}

type baseLogger struct {
	worker    *glg.Glg
	file      *os.File
	closeOnce sync.Once
}

func (l *baseLogger) Debug(format string, v ...interface{}) {
	l.worker.Debugf(format, v...)
}

func (l *baseLogger) Close() (err error) {
	l.closeOnce.Do(func() {
		if l.file != nil {
			err = l.file.Close()
		}
	})
	return
}

func New(logLevel LogLevel, outputFlags Output, logPath string) (Logger, error) {
	worker := glg.New().SetLevel(glg.LEVEL(logLevel))
	var fileWriter *os.File

	if isFlagOn(FileOut, outputFlags) {
		fileWriter = glg.FileWriter(logPath, 0644)
		worker = worker.AddWriter(fileWriter)
	}

	if isFlagOn(StdOut, outputFlags) && isFlagOn(FileOut, outputFlags) {
		worker = worker.SetMode(glg.BOTH)
	} else if isFlagOn(FileOut, outputFlags) {
		worker = worker.SetMode(glg.WRITER)
	} else if isFlagOn(StdOut, outputFlags) {
		worker = worker.SetMode(glg.STD)
	} else {
		worker = worker.SetMode(glg.NONE)
	}

	baseLogger := &baseLogger{
		worker: worker,
		file:   fileWriter,
	}

	return baseLogger, nil
}

func isFlagOn(singleFlag, combinedFlag Output) bool {
	return singleFlag&combinedFlag != 0
}
