package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger interface {
	Close() <-chan (bool)
	Info(format string, v ...interface{})
}

type defaultLogger struct {
	writer io.Writer
	file   *os.File
	queue  chan (string)
	done   chan (bool)
}

type DebugLogger struct {
	defaultLogger
}

func (l *defaultLogger) addQueue(prefix, format string, v ...interface{}) {
	l.queue <- prefix + fmt.Sprintf(format, v...)
}

func (l *defaultLogger) runLoggerDaemon() {
	for msg := range l.queue {
		l.writer.Write([]byte(time.Now().Format("2006-01-02 15:04:05 ")))
		l.writer.Write([]byte(msg + "\n"))
	}
	l.done <- true
}

func (l *DebugLogger) Info(format string, v ...interface{}) {
	if len(format) == 0 {
		return
	}
	go l.addQueue("[INFO]", format, v...)
}

func (l *DebugLogger) Close() <-chan (bool) {
	close(l.queue)
	<-l.done
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
		defaultLogger: defaultLogger{
			writer: io.MultiWriter(writers[:writersCount]...),
			file:   fileWriter,
			queue:  make(chan string, 20),
			done:   make(chan bool),
		},
	}

	go debugLogger.runLoggerDaemon()
	return debugLogger, nil
}

func isFlagOn(singleFlag, combinedFlag int) bool {
	return singleFlag&combinedFlag != 0
}
