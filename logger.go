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
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	// Fatal log a message and exit program
	Fatal(format string, v ...interface{})
}

type baseLogger struct {
	worker    *glg.Glg
	file      *os.File
	closeOnce sync.Once
}

func (l *baseLogger) Debug(format string, v ...interface{}) {
	l.worker.Debugf(format, v...)
}

func (l *baseLogger) Info(format string, v ...interface{}) {
	l.worker.Infof(format, v...)
}

func (l *baseLogger) Warn(format string, v ...interface{}) {
	l.worker.Warnf(format, v...)
}

func (l *baseLogger) Error(format string, v ...interface{}) {
	l.worker.Errorf(format, v...)
}

func (l *baseLogger) Fatal(format string, v ...interface{}) {
	l.worker.Fatalf(format, v...)
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
	worker := glg.New()
	var fileWriter *os.File

	if isFlagOn(StdOut, outputFlags) && isFlagOn(FileOut, outputFlags) {
		worker = worker.SetMode(glg.BOTH)
	} else if isFlagOn(FileOut, outputFlags) {
		worker = worker.SetMode(glg.WRITER)
	} else if isFlagOn(StdOut, outputFlags) {
		worker = worker.SetMode(glg.STD)
	} else {
		worker = worker.SetMode(glg.NONE)
	}

	if isFlagOn(FileOut, outputFlags) {
		fileWriter = glg.FileWriter(logPath, 0644)
		worker = worker.AddWriter(fileWriter)
	}

	baseLogger := &baseLogger{
		worker: worker.SetLevel(glg.LEVEL(logLevel)),
		file:   fileWriter,
	}

	return baseLogger, nil
}

func isFlagOn(singleFlag, combinedFlag Output) bool {
	return singleFlag&combinedFlag != 0
}
