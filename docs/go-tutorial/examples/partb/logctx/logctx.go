// Package logctx 把请求号放进 context，并在打日志时带上。
package logctx

import (
	"context"
	"log/slog"
)

type requestIDKey struct{}

// WithRequestID 返回带请求号的子 context。不要把可选业务参数塞进 context。
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID 读取请求号。没有时返回空字符串。
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// WithLogger 给 logger 加上从 ctx 里能拿到的请求号。
func WithLogger(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if id := RequestID(ctx); id != "" {
		return logger.With("request_id", id)
	}
	return logger
}
