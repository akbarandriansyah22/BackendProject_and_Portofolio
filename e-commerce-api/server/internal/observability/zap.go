package observability

import (
	"go.uber.org/zap"
)

// ZapLogger is an implementation of Logger using go.uber.org/zap
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger instance
func NewZapLogger() Logger {
	l, _ := zap.NewProduction()
	return &ZapLogger{logger: l}
}

// NewZapLoggerFrom wraps an existing zap logger.
func NewZapLoggerFrom(l *zap.Logger) Logger {
	return &ZapLogger{logger: l}
}

// Zap returns the underlying zap logger for repositories that require it.
func (zl *ZapLogger) Zap() *zap.Logger {
	return zl.logger
}

// Log logs a message at the specified level
func (zl *ZapLogger) Log(level Level, msg string, args ...interface{}) {
	switch level {
	case Debug:
		zl.logger.Sugar().Debugf(msg, args...)
	case Info:
		zl.logger.Sugar().Infof(msg, args...)
	case Warn:
		zl.logger.Sugar().Warnf(msg, args...)
	case Error:
		zl.logger.Sugar().Errorf(msg, args...)
	default:
		zl.logger.Sugar().Infof(msg, args...)
	}
}

// Debug logs a debug-level message
func (zl *ZapLogger) Debug(msg string, args ...interface{}) {
	zl.logger.Sugar().Debugf(msg, args...)
}

// Info logs an info-level message
func (zl *ZapLogger) Info(msg string, args ...interface{}) {
	zl.logger.Sugar().Infof(msg, args...)
}

// Warn logs a warn-level message
func (zl *ZapLogger) Warn(msg string, args ...interface{}) {
	zl.logger.Sugar().Warnf(msg, args...)
}

// Error logs an error-level message
func (zl *ZapLogger) Error(msg string, args ...interface{}) {
	zl.logger.Sugar().Errorf(msg, args...)
}

// Fatal logs a fatal-level message and exits
func (zl *ZapLogger) Fatal(msg string, args ...interface{}) {
	zl.logger.Sugar().Fatalf(msg, args...)
}

// Sync flushes the logger
func (zl *ZapLogger) Sync() error {
	return zl.logger.Sync()
}
