package logger_test

import (
	"fmt"
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
	logPath := "./testdata/newLoggerTest.log"
	defaultLogger := newLogger(t, logger.Debug, logger.FileOut, logPath)
	defer func() {
		<-defaultLogger.Close()
	}()

	logFile, err := os.Open(logPath)
	if os.IsNotExist(err) {
		t.Errorf("fail to create log %w", err)
	}

	logFile.Close()
	remove(t, logPath)
}

func TestInfo(t *testing.T) {
	logPath := "./testdata/infoTest.log"
	debugLogger := newLogger(t, logger.Debug, logger.FileOut, logPath)
	defer func() {
		<-debugLogger.Close()
	}()

	testCases := []struct {
		input, expect string
	}{
		{input: "Test", expect: "[INFO]Test"},
		{input: "Testwhitespace", expect: "[INFO]Testwhitespace"},
		{input: "", expect: ""},
		{},
	}

	for _, testCase := range testCases {
		debugLogger.Info(testCase.input)
	}

	logFile, err := os.Open(logPath)
	if err == nil {
		defer logFile.Close()
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("%dst case in Info Test", idx+1), func(t *testing.T) {
			var message string
			var date, time string
			fmt.Fscanf(logFile, "%s %s %s", &date, &time, &message)
			if message != testCase.expect {
				t.Errorf("input: %s, expect: %s, actual: %s", testCase.input, testCase.expect, message)
			}
		})
	}

	remove(t, logPath)
}

func TestReceiveMessageWhenLoggerIsClosed(t *testing.T) {
	defaultLogger := newLogger(t, logger.Debug, logger.Discard, "")
	<-defaultLogger.Close()
}

func remove(t *testing.T, filePath string) {
	t.Helper()
	t.Cleanup(func() {
		err := os.Remove(filePath)
		if err != nil {
			t.Error(err)
		}
	})
}

func newLogger(t *testing.T, logLevel, flags int, logPath string) logger.Logger {
	t.Helper()

	logger, err := logger.New(logLevel, flags, logPath)

	if err != nil {
		t.FailNow()
	}

	return logger
}
