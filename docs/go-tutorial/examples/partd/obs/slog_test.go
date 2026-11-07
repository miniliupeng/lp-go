package obs

import (
	"bytes"
	"log/slog"
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
}
