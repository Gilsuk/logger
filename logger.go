package logger

import (
	"io"
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

type DebugLogger struct {
	baseLogger
}

func (l *DebugLogger) Debug(format string, v ...interface{}) {
	l.worker.Debugf(format, v...)
}

func (l *DebugLogger) Close() (err error) {
	l.closeOnce.Do(func() {
		if l.file != nil {
			err = l.file.Close()
		}
	})
	return
}

func New(logLevel LogLevel, outputFlags Output, logPath string) (Logger, error) {
	writers := make([]io.Writer, 0)
	var fileWriter *os.File
	writersCount := 0

	if isFlagOn(FileOut, outputFlags) {
		var err error
		fileWriter, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writersCount++
		writers = append(writers, fileWriter)
	}

	if isFlagOn(StdOut, outputFlags) {
		writersCount++
		writers = append(writers, os.Stdout)
	}

	if isFlagOn(Discard, outputFlags) {
		writersCount++
		writers = append(writers, io.Discard)
	}

	debugLogger := &DebugLogger{
		baseLogger: baseLogger{
			worker: glg.New().AddWriter(io.MultiWriter(writers[:writersCount]...)).SetLevel(glg.DEBG).SetMode(glg.WRITER),
			file:   fileWriter,
		},
	}

	return debugLogger, nil
}

func isFlagOn(singleFlag, combinedFlag Output) bool {
	return singleFlag&combinedFlag != 0
}
