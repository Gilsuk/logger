package logger

type Logger interface {
}

type DefaultLogger struct {
}

func New(outputFlags int, logPath string) (Logger, error) {
	return &DefaultLogger{}, nil
}
