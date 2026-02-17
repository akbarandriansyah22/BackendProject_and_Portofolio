package observability

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type Level string

const (
	Debug Level = "DEBUG"
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
)

type Logger interface {
	Log(level Level, msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Sync() error
}

type standardLogger struct {
	logger *log.Logger
}

func NewLogger() Logger {
	return &standardLogger{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *standardLogger) Log(level Level, msg string, args ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)

	_, file, line, ok := runtime.Caller(3)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	message := fmt.Sprintf(msg, args...)
	l.logger.Printf("[%s] [%s] [%s] %s", timestamp, level, caller, message)
}

func (l *standardLogger) Debug(msg string, args ...interface{}) {
	l.Log(Debug, msg, args...)
}

func (l *standardLogger) Info(msg string, args ...interface{}) {
	l.Log(Info, msg, args...)
}

func (l *standardLogger) Warn(msg string, args ...interface{}) {
	l.Log(Warn, msg, args...)
}

func (l *standardLogger) Error(msg string, args ...interface{}) {
	l.Log(Error, msg, args...)
}

func (l *standardLogger) Fatal(msg string, args ...interface{}) {
	l.Log(Error, msg, args...)
	os.Exit(1)
}

func (l *standardLogger) Sync() error {
	return nil
}
