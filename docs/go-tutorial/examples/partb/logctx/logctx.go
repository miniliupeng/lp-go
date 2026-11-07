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

// RunSteps 做最多 steps 次工作。ctx 已结束时立刻返回已完成的步数。
// 返回的错误是 context.Cause(ctx)：WithCancelCause 取消时是你传入的原因，
// 普通取消时是 context.Canceled，超时是 context.DeadlineExceeded。
// 未取消时 Cause 返回 nil，这里也返回 nil。
//
// 先看 ctx.Err()，不要用 select 的 default 和 ctx.Done() 抢。两者都就绪时 select 会随机选一个。
func RunSteps(ctx context.Context, steps int) (int, error) {
	for n := 0; n < steps; n++ {
		if err := ctx.Err(); err != nil {
			return n, context.Cause(ctx)
		}
	}
	return steps, nil
}

// WithLogger 给 logger 加上从 ctx 里能拿到的请求号。
func WithLogger(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if id := RequestID(ctx); id != "" {
		return logger.With("request_id", id)
	}
	return logger
}
