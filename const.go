package logger

import "github.com/kpango/glg"

type Output int
type LogLevel glg.LEVEL

const (
	StdOut Output = 1 << iota
	FileOut
	Discard
)

const (
	Debug LogLevel = LogLevel(glg.DEBG)
	Info           = LogLevel(glg.INFO)
	Warn           = LogLevel(glg.WARN)
	Error          = LogLevel(glg.ERR)
	Fatal          = LogLevel(glg.FATAL)
)
