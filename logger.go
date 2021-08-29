package logger

type Logger interface {
}

type DefaultLogger struct {
}

func New(logLevel, outputFlags int, logPath string) (Logger, error) {
	return &DefaultLogger{}, nil
}
