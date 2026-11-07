package logctx

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"testing/synctest"
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

func TestRunStepsStopsWhenCanceled(t *testing.T) {
	cause := errors.New("upstream")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)
	n, err := RunSteps(ctx, 10)
	if n != 0 || !errors.Is(err, cause) {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatal(ctx.Err())
	}
}

func TestRunStepsTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		synctest.Wait()
		if ctx.Err() != nil {
			t.Fatal(ctx.Err())
		}
		time.Sleep(time.Second)
		synctest.Wait()
		n, err := RunSteps(ctx, 3)
		if n != 0 || !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("n=%d err=%v", n, err)
		}
	})
}
