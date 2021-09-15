package logger_test

import (
	"testing"

	"github.com/gilsuk/logger"
)

func TestOutFlags(t *testing.T) {
	if !isFlagOn(logger.Discard, logger.Discard|logger.FileOut) {
		t.Fail()
	}

	if isFlagOn(logger.Discard, logger.StdOut|logger.FileOut) {
		t.Fail()
	}
}

func isFlagOn(singleFlag, combinedFlag logger.Output) bool {
	return singleFlag&combinedFlag != 0
}
