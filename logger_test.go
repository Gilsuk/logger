package logger_test

import (
	"testing"

	"github.com/gilsuk/logger"
)

func TestNewLogger(t *testing.T) {
	_, err := logger.New(logger.Discard, "")

	if err != nil {
		t.Fail()
	}
}
