package httplog

import (
	"context"
)

// Logger represents an interface of a logger
type Logger interface {
	Debugf(msg string, args ...any)
	Warnf(msg string, args ...any)
	Infof(msg string, args ...any)
	Errorf(msg string, args ...any)
	Panicf(msg string, args ...any)

	Debugw(fields []Field, msg string, args ...any)
	Warnw(fields []Field, msg string, args ...any)
	Infow(fields []Field, msg string, args ...any)
	Errorw(fields []Field, msg string, args ...any)
	Panicw(fields []Field, msg string, args ...any)

	Debugcf(ctx context.Context, msg string, args ...any)
	Warncf(ctx context.Context, msg string, args ...any)
	Infocf(ctx context.Context, msg string, args ...any)
	Errorcf(ctx context.Context, msg string, args ...any)
	Paniccf(ctx context.Context, msg string, args ...any)

	Debugcw(ctx context.Context, fields []Field, msg string, args ...any)
	Warncw(ctx context.Context, fields []Field, msg string, args ...any)
	Infocw(ctx context.Context, fields []Field, msg string, args ...any)
	Errorcw(ctx context.Context, fields []Field, msg string, args ...any)
	Paniccw(ctx context.Context, fields []Field, msg string, args ...any)

	Named(name string) Logger
}

var globalLogger Logger

// SetGlobalLogger sets the global logger
func SetGlobalLogger(l Logger) {
	globalLogger = l
	l.Infof("Global logger is set")
}

// GetLogger returns the global logger
func GetLogger() Logger {
	if globalLogger == nil {
		panic("global logger is not set yet, please call SetGlobalLogger() first")
	}
	return globalLogger
}
