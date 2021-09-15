package logger_test

import (
	"fmt"
	"os"
	"testing"
	"time"

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
		t.Errorf("fail to create log %v", err)
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
		format, expect string
		vars           []interface{}
	}{
		{format: "Test", expect: "[INFO]Test"},
		{format: "TestFormat%sNumber%d",
			vars:   []interface{}{"String", 10},
			expect: "[INFO]TestFormatStringNumber10"},
		{format: "", expect: ""},
		{},
	}

	for _, testCase := range testCases {
		debugLogger.Info(testCase.format, testCase.vars...)
	}

	logFile, err := os.Open(logPath)
	if err == nil {
		defer logFile.Close()
	} else {
		t.Errorf("Can not open file")
	}

	time.Sleep(time.Second)
	done := make(chan (bool), 1)
	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("%dst case in Info Test", idx+1), func(t *testing.T) {
			var message string
			var date, time string
			fmt.Fscanf(logFile, "%s %s %s", &date, &time, &message)
			if message != testCase.expect {
				t.Errorf("input: %s, expect: %s, actual: %s", testCase.format, testCase.expect, message)
			}
			if idx+1 == len(testCases) {
				done <- true
			}
		})
	}

	<-done
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
