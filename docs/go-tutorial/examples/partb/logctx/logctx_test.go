package logctx

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestRequestIDAndLog(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx = WithRequestID(ctx, "req-7")
	if RequestID(ctx) != "req-7" {
		t.Fatal(RequestID(ctx))
	}
	if RequestID(context.Background()) != "" {
		t.Fatal("background")
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	WithLogger(ctx, logger).Info("hello")
	if !strings.Contains(buf.String(), `"request_id":"req-7"`) {
		t.Fatalf("log = %s", buf.String())
	}
	if err := ctx.Err(); err != nil {
		t.Fatal(err)
	}
}
