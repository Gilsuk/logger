package logger

import (
	"io"
	"os"
)

type Logger interface {
	Close() <-chan (bool)
}

type defaultLogger struct {
	writer io.Writer
	file   *os.File
}

type DebugLogger struct {
	defaultLogger
}

func (l *DebugLogger) Close() <-chan (bool) {
	doneChan := make(chan bool)
	go func() {
		if l.file != nil {
			l.file.Close()
		}
		doneChan <- true
	}()
	return doneChan
}

func New(logLevel, outputFlags int, logPath string) (Logger, error) {
	writers := make([]io.Writer, 3)
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

	return &DebugLogger{
		defaultLogger: defaultLogger{
			writer: io.MultiWriter(writers[:writersCount]...),
			file:   fileWriter,
		},
	}, nil
}

func isFlagOn(singleFlag, combinedFlag int) bool {
	return singleFlag&combinedFlag != 0
}
