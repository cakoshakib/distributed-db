package commons

import (
	"context"
	"go.uber.org/zap"
)

type contextKey string

const LoggerKey = contextKey("logger")

func LoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}
	return logger
}
