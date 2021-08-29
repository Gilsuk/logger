package logger_test

import (
	"testing"

	"github.com/gilsuk/logger"
)

func TestNewLogger(t *testing.T) {
	defaultLogger := newLogger(t, logger.Debug, logger.Discard, "")
	if _, ok := defaultLogger.(*logger.DefaultLogger); !ok {
		t.Fail()
	}
}

func newLogger(t *testing.T, logLevel, flags int, logPath string) logger.Logger {
	t.Helper()

	logger, err := logger.New(logLevel, flags, logPath)

	if err != nil {
		t.Fail()
	}

	return logger
}
