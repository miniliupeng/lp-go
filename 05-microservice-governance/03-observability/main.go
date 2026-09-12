package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 1. OpenTelemetry W3C Trace Context 规范与链路贯穿
// ============================================================================

type traceKey struct{}

// SpanContext 链路追踪上下文
type SpanContext struct {
	TraceID string // 16 字节 (32 hex 字符)
	SpanID  string // 8 字节 (16 hex 字符)
	Sampled bool
}

// GenerateRandomHex 生成指定长度的随机十六进制串
func GenerateRandomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// NewTraceContext 初始化新根链路
func NewTraceContext() *SpanContext {
	return &SpanContext{
		TraceID: GenerateRandomHex(16),
		SpanID:  GenerateRandomHex(8),
		Sampled: true,
	}
}

// ChildSpan 基于父上下文派生子 Span
func (sc *SpanContext) ChildSpan() *SpanContext {
	return &SpanContext{
		TraceID: sc.TraceID, // 保持 TraceID 不变贯穿全局
		SpanID:  GenerateRandomHex(8),
		Sampled: sc.Sampled,
	}
}

// ContextWithSpan 将 SpanContext 注入 Go Context
func ContextWithSpan(ctx context.Context, span *SpanContext) context.Context {
	return context.WithValue(ctx, traceKey{}, span)
}

// SpanFromContext 从 Context 提取 SpanContext
func SpanFromContext(ctx context.Context) *SpanContext {
	if val := ctx.Value(traceKey{}); val != nil {
		if span, ok := val.(*SpanContext); ok {
			return span
		}
	}
	return nil
}

// InjectHTTP 将链路信息注入 W3C 标准 traceparent 请求头
// 格式: 00-{trace_id}-{span_id}-{flags}
func (sc *SpanContext) InjectHTTP(header http.Header) {
	flags := "00"
	if sc.Sampled {
		flags = "01"
	}
	header.Set("traceparent", fmt.Sprintf("00-%s-%s-%s", sc.TraceID, sc.SpanID, flags))
}

// ExtractHTTP 从 HTTP Header 解析 traceparent
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

// ============================================================================
// 2. Prometheus RED 指标模型 (Rate, Errors, Duration)
// ============================================================================

// REDMetrics 统计接口吞吐率、错误率及分位耗时
type REDMetrics struct {
	TotalRequests  atomic.Uint64
	TotalErrors    atomic.Uint64
	DurationBuckets [5]atomic.Uint64 // 耗时桶：<10ms, <50ms, <100ms, <500ms, >=500ms
	mu             sync.RWMutex
	samples        []float64
}

func NewREDMetrics() *REDMetrics {
	return &REDMetrics{
		samples: make([]float64, 0, 1000),
	}
}

func (m *REDMetrics) Record(duration time.Duration, isErr bool) {
	m.TotalRequests.Add(1)
	if isErr {
		m.TotalErrors.Add(1)
	}

	ms := float64(duration.Microseconds()) / 1000.0
	m.mu.Lock()
	if len(m.samples) < 10000 {
		m.samples = append(m.samples, ms)
	}
	m.mu.Unlock()

	switch {
	case ms < 10:
		m.DurationBuckets[0].Add(1)
	case ms < 50:
		m.DurationBuckets[1].Add(1)
	case ms < 100:
		m.DurationBuckets[2].Add(1)
	case ms < 500:
		m.DurationBuckets[3].Add(1)
	default:
		m.DurationBuckets[4].Add(1)
	}
}

// ErrorRate 返回错误率百分比
func (m *REDMetrics) ErrorRate() float64 {
	total := m.TotalRequests.Load()
	if total == 0 {
		return 0.0
	}
	return (float64(m.TotalErrors.Load()) / float64(total)) * 100
}

// P99Duration 计算近似 P99 耗时 (毫秒)
func (m *REDMetrics) P99Duration() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := len(m.samples)
	if n == 0 {
		return 0
	}
	idx := int(math.Floor(0.99 * float64(n)))
	if idx >= n {
		idx = n - 1
	}
	return m.samples[idx]
}

func main() {
	fmt.Println("=== 阶段三 专题03：可观测性（OpenTelemetry 链路贯穿 + Prometheus RED 模型） ===")

	// 1. 初始化链路并模拟跨服务 RPC 调用透传
	rootSpan := NewTraceContext()
	fmt.Printf("[Gateway 发起请求] TraceID=%s, SpanID=%s\n", rootSpan.TraceID, rootSpan.SpanID)

	// 模拟网关向下游订单服务发送 HTTP 请求
	req, _ := http.NewRequest("GET", "/order/create", nil)
	rootSpan.InjectHTTP(req.Header)

	// 下游服务解析 HTTP Header 并派生子 Span
	orderSpan := ExtractHTTP(req.Header).ChildSpan()
	ctx := ContextWithSpan(context.Background(), orderSpan)

	// 从 Context 取出 Trace 信息打日志
	extracted := SpanFromContext(ctx)
	fmt.Printf("[OrderService 接收请求] 链路贯穿 TraceID=%s, 子SpanID=%s (父SpanID=%s)\n",
		extracted.TraceID, extracted.SpanID, rootSpan.SpanID)

	// 2. RED 指标打点
	metrics := NewREDMetrics()
	for i := 0; i < 100; i++ {
		isErr := (i%10 == 0) // 10% 错误率
		dur := time.Duration(i*2) * time.Millisecond
		metrics.Record(dur, isErr)
	}

	fmt.Printf("\n[RED 模型指标汇总]\n")
	fmt.Printf("总请求数 (Rate): %d\n", metrics.TotalRequests.Load())
	fmt.Printf("错误率 (Errors): %.2f%%\n", metrics.ErrorRate())
	fmt.Printf("P99 耗时 (Duration): %.2f ms\n", metrics.P99Duration())
}
