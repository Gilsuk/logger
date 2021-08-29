package logger

type Logger interface {
	Close() <-chan (bool)
}

type DefaultLogger struct {
}

func (*DefaultLogger) Close() <-chan (bool) {
	doneChan := make(chan bool)
	go func() {
		doneChan <- true
	}()
	return doneChan
}

func New(logLevel, outputFlags int, logPath string) (Logger, error) {
	return &DefaultLogger{}, nil
}
