package main

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestTraceContextPropagation(t *testing.T) {
	root := NewTraceContext()
	req, _ := http.NewRequest("GET", "/api/v1/user", nil)
	root.InjectHTTP(req.Header)

	extracted := ExtractHTTP(req.Header)
	if extracted == nil {
		t.Fatal("failed to extract traceparent")
	}

	if extracted.TraceID != root.TraceID {
		t.Fatalf("traceID mismatch: %s vs %s", extracted.TraceID, root.TraceID)
	}

	child := extracted.ChildSpan()
	if child.TraceID != root.TraceID {
		t.Fatalf("child traceID mismatch: %s vs %s", child.TraceID, root.TraceID)
	}
	if child.SpanID == root.SpanID {
		t.Fatal("child SpanID should be different from parent")
	}

	ctx := ContextWithSpan(context.Background(), child)
	fromCtx := SpanFromContext(ctx)
	if fromCtx == nil || fromCtx.SpanID != child.SpanID {
		t.Fatal("failed to get span from ctx")
	}
}

func TestREDMetrics(t *testing.T) {
	metrics := NewREDMetrics()

	metrics.Record(5*time.Millisecond, false)
	metrics.Record(15*time.Millisecond, false)
	metrics.Record(60*time.Millisecond, true)
	metrics.Record(120*time.Millisecond, false)

	if metrics.TotalRequests.Load() != 4 {
		t.Fatalf("expected 4 requests, got %d", metrics.TotalRequests.Load())
	}
	if metrics.TotalErrors.Load() != 1 {
		t.Fatalf("expected 1 error, got %d", metrics.TotalErrors.Load())
	}
	if metrics.ErrorRate() != 25.0 {
		t.Fatalf("expected 25%% error rate, got %.2f", metrics.ErrorRate())
	}
}

func BenchmarkTraceInjectExtract(b *testing.B) {
	root := NewTraceContext()
	header := make(http.Header)
	b.ResetTimer()
	for b.Loop() {
		root.InjectHTTP(header)
		_ = ExtractHTTP(header)
	}
}

func BenchmarkREDMetricsRecord(b *testing.B) {
	metrics := NewREDMetrics()
	dur := 10 * time.Millisecond
	b.ResetTimer()
	for b.Loop() {
		metrics.Record(dur, false)
	}
}
