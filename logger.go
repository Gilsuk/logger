package logger

type Logger interface {
	Close() <-chan (bool)
}

type DebugLogger struct {
}

func (*DebugLogger) Close() <-chan (bool) {
	doneChan := make(chan bool)
	go func() {
		doneChan <- true
	}()
	return doneChan
}

func New(logLevel, outputFlags int, logPath string) (Logger, error) {
	return &DebugLogger{}, nil
}
