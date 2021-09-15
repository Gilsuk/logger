package logger_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gilsuk/logger"
)

func TestNewFileLogger(t *testing.T) {
	logPath := "./testdata/newLoggerTest.log"
	baseLogger := newLogger(t, logger.Debug, logger.FileOut, logPath)
	defer baseLogger.Close()

	logFile, err := os.Open(logPath)
	if os.IsNotExist(err) {
		t.Errorf("fail to create log %v", err)
	}

	logFile.Close()
	remove(t, logPath)
}

func TestDebug(t *testing.T) {
	logPath := "./testdata/infoTest.log"
	baseLogger := newLogger(t, logger.Debug, logger.FileOut, logPath)
	defer baseLogger.Close()

	testCases := []struct {
		format, expect string
		vars           []interface{}
	}{
		{format: "Test", expect: "Test"},
		{format: "TestFormat%sNumber%d",
			vars:   []interface{}{"String", 10},
			expect: "TestFormatStringNumber10"},
		{format: "", expect: ""},
		{},
	}

	for _, testCase := range testCases {
		baseLogger.Debug(testCase.format, testCase.vars...)
	}

	logFile, err := os.Open(logPath)
	if err == nil {
		defer logFile.Close()
	} else {
		t.Errorf("Can not open file")
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("%dst case in Debug Test", idx+1), func(t *testing.T) {
			var logLevel, message string
			var date, time string
			fmt.Fscanf(logFile, "%s %s %s %s", &date, &time, &logLevel, &message)
			if message != testCase.expect {
				t.Errorf("input: %s, expect: %s, actual: %s", testCase.format, testCase.expect, message)
			}
		})
	}

	remove(t, logPath)
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

func newLogger(t *testing.T, logLevel logger.LogLevel, flags logger.Output, logPath string) logger.Logger {
	t.Helper()

	logger, err := logger.New(logLevel, flags, logPath)

	if err != nil {
		t.FailNow()
	}

	return logger
}
