package logger

import (
	"go.uber.org/zap"
)

// Logger wraps a Zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new Zap logger
func NewLogger() (*Logger, error) {
	// Use production config for structured logging
	cfg := zap.NewProductionConfig()
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.Logger.Debug(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.Logger.Info(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.Logger.Error(msg)
}
