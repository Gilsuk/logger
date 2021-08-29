package logger_test

import (
	"testing"

	"github.com/gilsuk/logger"
)

func TestNewLogger(t *testing.T) {
	debugLogger := newLogger(t, logger.Debug, logger.Discard, "")
	if _, ok := debugLogger.(*logger.DebugLogger); !ok {
		t.Fail()
	}
	debugLogger.Close()
}

func TestReceiveMessageWhenLoggerIsClosed(t *testing.T) {
	defaultLogger := newLogger(t, logger.Debug, logger.Discard, "")
	<-defaultLogger.Close()
}

func newLogger(t *testing.T, logLevel, flags int, logPath string) logger.Logger {
	t.Helper()

	logger, err := logger.New(logLevel, flags, logPath)

	if err != nil {
		t.Fail()
	}

	return logger
}
