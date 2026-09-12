package main

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSSEGatewayHandler(t *testing.T) {
	handler := NewSSEGatewayHandler(10, 100*time.Millisecond)

	req := httptest.NewRequest("GET", "/api/chat/completions", nil)
	w := httptest.NewRecorder()

	handler.HandleChatStream(w, req)

	resp := w.Result()
	if resp.Header.Get("Content-Type") != "text/event-stream" {
		t.Fatalf("expected text/event-stream, got %s", resp.Header.Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "event: message") {
		t.Fatalf("expected sse events, got: %s", body)
	}
	if !strings.Contains(body, "Go") {
		t.Fatalf("expected generated token 'Go', got: %s", body)
	}
}

func TestBoundedStreamChannelBackpressure(t *testing.T) {
	stream := NewBoundedStreamChannel(1)
	ctx := context.Background()

	err1 := stream.Push(ctx, &StreamEvent{ID: 1}, 10*time.Millisecond)
	if err1 != nil {
		t.Fatalf("first push should succeed, got: %v", err1)
	}

	// 触发反压超时
	err2 := stream.Push(ctx, &StreamEvent{ID: 2}, 20*time.Millisecond)
	if err2 == nil {
		t.Fatal("second push should trigger backpressure error")
	}

	stream.Close()
}

func BenchmarkFormatSSE(b *testing.B) {
	evt := &StreamEvent{
		ID:        100,
		DeltaText: "Artificial Intelligence",
		Finished:  false,
	}

	b.ResetTimer()
	for b.Loop() {
		_ = FormatSSE(evt)
	}
}
