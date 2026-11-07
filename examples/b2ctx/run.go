// Package logctx 把请求号放进 context，并在循环里观察取消。
package logctx

import (
	"context"
	"log/slog"
)

type requestIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// RunSteps 最多做 steps 步。ctx 已结束时返回已完成的步数和 context.Cause。
// 先看 ctx.Err()。ctx 已取消时，select 的 Done 和 default 都就绪，会随机选一个。
func RunSteps(ctx context.Context, steps int) (int, error) {
	for n := 0; n < steps; n++ {
		if err := ctx.Err(); err != nil {
			return n, context.Cause(ctx)
		}
	}
	return steps, nil
}

func WithLogger(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if id := RequestID(ctx); id != "" {
		return logger.With("request_id", id)
	}
	return logger
}
