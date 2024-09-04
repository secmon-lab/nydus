package logging

import (
	"context"
	"log/slog"
	"sync"
)

var (
	defaultMutex  sync.RWMutex
	defaultLogger = slog.Default()
)

type ctxKey struct{}

func Inject(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func From(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return logger
	}

	return Default()
}

func SetDefault(logger *slog.Logger) {
	defaultMutex.Lock()
	defer defaultMutex.Unlock()
	defaultLogger = logger
}

func Default() *slog.Logger {
	defaultMutex.RLock()
	defer defaultMutex.RUnlock()
	return defaultLogger
}
