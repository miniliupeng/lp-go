package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type traceCtxKey struct{}

// SpanContext 链路追踪上下文契约 (遵循 W3C Trace Context 规范)
type SpanContext struct {
	TraceID string // 16 字节十六进制字符 (32 hex chars)
	SpanID  string // 8 字节十六进制字符 (16 hex chars)
	Sampled bool
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// NewRootSpan 创建全新的根链路
func NewRootSpan() *SpanContext {
	return &SpanContext{
		TraceID: randomHex(16),
		SpanID:  randomHex(8),
		Sampled: true,
	}
}

// ChildSpan 基于当前父上下文派生子 Span (维持全局相同 TraceID)
func (s *SpanContext) ChildSpan() *SpanContext {
	return &SpanContext{
		TraceID: s.TraceID,
		SpanID:  randomHex(8),
		Sampled: s.Sampled,
	}
}

// WithContext 将 SpanContext 注入 Go Context
func WithContext(ctx context.Context, span *SpanContext) context.Context {
	return context.WithValue(ctx, traceCtxKey{}, span)
}

// FromContext 从 Context 安全提取 SpanContext
func FromContext(ctx context.Context) *SpanContext {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(traceCtxKey{}); v != nil {
		if span, ok := v.(*SpanContext); ok {
			return span
		}
	}
	return nil
}

// ExtractHTTP 从 HTTP Header 中解析 W3C traceparent (格式: 00-{trace_id}-{span_id}-{flags})
func ExtractHTTP(header http.Header) *SpanContext {
	tp := header.Get("traceparent")
	if len(tp) < 55 {
		return nil
	}
	var version, traceID, spanID, flags string
	_, err := fmt.Sscanf(tp, "%2s-%32s-%16s-%2s", &version, &traceID, &spanID, &flags)
	if err != nil || version != "00" {
		return nil
	}
	return &SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
		Sampled: flags == "01",
	}
}

// InjectHTTP 将 SpanContext 注入 HTTP 请求头
func (s *SpanContext) InjectHTTP(header http.Header) {
	flags := "00"
	if s.Sampled {
		flags = "01"
	}
	header.Set("traceparent", fmt.Sprintf("00-%s-%s-%s", s.TraceID, s.SpanID, flags))
}
