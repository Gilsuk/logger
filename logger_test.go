package logger_test

import (
	"testing"

	"github.com/gilsuk/logger"
)

func TestNewLogger(t *testing.T) {
	newLogger(t, logger.Discard, "")
}

func newLogger(t *testing.T, flags int, logPath string) logger.Logger {
	t.Helper()

	logger, err := logger.New(logger.Discard, "")

	if err != nil {
		t.Fail()
	}

	return logger
}
