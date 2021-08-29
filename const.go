package logger

const (
	StdOut int = 1 << iota
	FileOut
	Discard
)

const (
	Debug = iota
	Info
	Warn
	Error
)
