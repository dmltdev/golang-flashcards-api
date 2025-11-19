package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	service string
}

var defaultLogger *Logger

func init() {
	defaultLogger = New("app")
}

func New(service string) *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	return &Logger {
		Logger: slog.New(handler),
		service: service,
	}
}

func Default() *Logger {
	return defaultLogger
}

func (l *Logger) Info(msg string, args ...any) {
    l.Logger.Info(msg, append([]any{"service", l.service}, args...)...)
}

func (l *Logger) Error(msg string, err error, args ...any) {
    allArgs := append([]any{"service", l.service, "error", err}, args...)
    l.Logger.Error(msg, allArgs...)
}

func (l *Logger) Debug(msg string, args ...any) {
    l.Logger.Debug(msg, append([]any{"service", l.service}, args...)...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, append([]any{"service", l.service}, args...)...)
}

func (l *Logger) WithFields(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		service: l.service,
	}
}