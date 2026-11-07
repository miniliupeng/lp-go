package svc

import (
	"bytes"
	"log/slog"
	"runtime/metrics"
	"strings"
	"testing"
	"time"
)

func TestBreaker(t *testing.T) {
	start := time.Unix(0, 0)
	b := NewBreaker(2, time.Second)
	if !b.Allow(start) {
		t.Fatal("closed")
	}
	b.Failure(start)
	if !b.Allow(start) {
		t.Fatal("still closed after one failure")
	}
	b.Failure(start)
	if b.Allow(start) {
		t.Fatal("should be open")
	}
	if b.Allow(start.Add(time.Second)) != true {
		t.Fatal("probe")
	}
	if b.Allow(start.Add(time.Second)) {
		t.Fatal("second probe")
	}
	b.Failure(start.Add(time.Second))
	if b.Allow(start.Add(time.Second)) {
		t.Fatal("reopened")
	}
	if !b.Allow(start.Add(2 * time.Second)) {
		t.Fatal("probe again")
	}
	b.Success()
	if !b.Allow(start.Add(2 * time.Second)) {
		t.Fatal("closed after success")
	}
}

func TestSlogAndMetrics(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("ready", "port", 8080)
	logger.Debug("hidden")
	text := buf.String()
	if !strings.Contains(text, `"msg":"ready"`) || !strings.Contains(text, `"port":8080`) || strings.Contains(text, "hidden") {
		t.Fatal(text)
	}
	if logger.Enabled(t.Context(), slog.LevelDebug) || !logger.Enabled(t.Context(), slog.LevelInfo) {
		t.Fatal("level")
	}
	samples := []metrics.Sample{
		{Name: "/memory/classes/total:bytes"},
		{Name: "/memory/classes/heap/released:bytes"},
	}
	metrics.Read(samples)
	if samples[0].Value.Kind() == metrics.KindBad || samples[0].Value.Uint64() == 0 {
		t.Fatal(samples[0].Value.Kind())
	}
	if samples[1].Value.Kind() == metrics.KindBad {
		t.Fatal("released")
	}
}
