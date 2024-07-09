package util

import (
	"errors"
	"log/slog"
)

// Ptr is a helper that returns a pointer of whatever is passed in
func Ptr[T any](t T) *T {
	return &t
}

// LogLevelInt convers the string version of a log level to a slog.Level
// Valid inputs: INFO, WARN, ERROR, DEBUG
func LogLevelInt(level string) (slog.Level, error) {
	switch level {
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	case "DEBUG":
		return slog.LevelDebug, nil
	default:
		return slog.Level(-1), errors.New("invalid log level. Please set one of: INFO, WARN, ERROR, DEBUG")
	}
}
