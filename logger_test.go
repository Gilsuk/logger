package logger_test

import (
	"os"
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

func TestNewFileLogger(t *testing.T) {
	defaultLogger := newLogger(t, logger.Debug, logger.FileOut, "./testdata/infoTest.log")
	defer func() {
		<-defaultLogger.Close()
	}()

	logFile, err := os.Open("./testdata/infoTest.log")
	if os.IsNotExist(err) {
		t.Errorf("fail to create log %w", err)
	}

	logFile.Close()
	t.Cleanup(func() {
		err := os.Remove("./testdata/infoTest.log")
		if err != nil {
			t.Error(err)
		}
	})

}

func TestReceiveMessageWhenLoggerIsClosed(t *testing.T) {
	defaultLogger := newLogger(t, logger.Debug, logger.Discard, "")
	<-defaultLogger.Close()
}

func newLogger(t *testing.T, logLevel, flags int, logPath string) logger.Logger {
	t.Helper()

	logger, err := logger.New(logLevel, flags, logPath)

	if err != nil {
		t.FailNow()
	}

	return logger
}
