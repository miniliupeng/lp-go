package obs

import (
	"bytes"
	"log/slog"
	"runtime/metrics"
	"strings"
	"testing"
)

func TestJSONLogHasLevelAndMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("ready", "port", 8080)
	logger.Debug("hidden")
	text := buf.String()
	if !strings.Contains(text, `"msg":"ready"`) || !strings.Contains(text, `"port":8080`) {
		t.Fatalf("log = %s", text)
	}
	if strings.Contains(text, "hidden") {
		t.Fatalf("debug leaked: %s", text)
	}
	if logger.Enabled(t.Context(), slog.LevelDebug) {
		t.Fatal("debug level is enabled")
	}
	if !logger.Enabled(t.Context(), slog.LevelInfo) {
		t.Fatal("info level is disabled")
	}
}

func TestRuntimeMetricsReadable(t *testing.T) {
	samples := []metrics.Sample{
		{Name: "/memory/classes/total:bytes"},
		{Name: "/memory/classes/heap/released:bytes"},
	}
	metrics.Read(samples)
	if samples[0].Value.Kind() == metrics.KindBad || samples[0].Value.Uint64() == 0 {
		t.Fatalf("total %#v", samples[0].Value)
	}
	if samples[1].Value.Kind() == metrics.KindBad {
		t.Fatalf("released %#v", samples[1].Value)
	}
}
