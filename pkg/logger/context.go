package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type contextKey int

const (
	loggerKey contextKey = iota
)

func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) zerolog.Logger {
	logger, ok := ctx.Value(loggerKey).(zerolog.Logger)
	if !ok {
		return zerolog.Nop()
	}

	return logger
}
