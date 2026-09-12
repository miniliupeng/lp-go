package logger

import (
	"context"
	"log/slog"
	"os"

	"lp-go/07-capstone-project/ai-gateway-service/pkg/trace"
)

// TraceHandler 自动拦截 slog 日志，将 Context 中的 TraceID 自动追加到日志属性中
type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if span := trace.FromContext(ctx); span != nil {
			r.AddAttrs(
				slog.String("trace_id", span.TraceID),
				slog.String("span_id", span.SpanID),
			)
		}
	}
	return h.Handler.Handle(ctx, r)
}

// InitLogger 初始化全局高性能结构化日志
func InitLogger() *slog.Logger {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	customHandler := &TraceHandler{Handler: baseHandler}
	l := slog.New(customHandler)
	slog.SetDefault(l)
	return l
}
