package ctxt

import (
	"context"

	"github.com/rs/zerolog"
)

var logCtxKey = contextKey("logCtx")

// LogContext stores the logger in the Context.
type LogContext struct {
	logger *zerolog.Logger
}

// InitLogContext returns a new context with a new LogContext.
func InitLogContext(ctx context.Context, log *zerolog.Logger) context.Context {
	logCtx := &LogContext{logger: log}
	return context.WithValue(ctx, logCtxKey, logCtx)
}

// SetLogger updates the logger in log context.
// Note that if the context is not initialized, the logger is not updated.
func SetLogger(ctx context.Context, log *zerolog.Logger) {
	if logCtx := getLogContext(ctx); logCtx != nil {
		logCtx.logger = log
	}
}

// GetLogger returns the zerolog logger from context. It returns nil if not available.
func GetLogger(ctx context.Context) *zerolog.Logger {
	if logCtx := getLogContext(ctx); logCtx != nil {
		return logCtx.logger
	}
	return nil
}

// getLogContext retrieves the logContext from standard context.
func getLogContext(ctx context.Context) *LogContext {
	if logCtx, ok := ctx.Value(logCtxKey).(*LogContext); ok {
		return logCtx
	}
	return nil
}
