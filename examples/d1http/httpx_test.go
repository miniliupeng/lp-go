package httpx

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestWildcardAndZeroTimeouts(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, r.PathValue("path"))
	})
	srv := &http.Server{Handler: mux}
	if srv.ReadTimeout != 0 || srv.ReadHeaderTimeout != 0 || srv.WriteTimeout != 0 || srv.IdleTimeout != 0 {
		t.Fatal("zero value should mean no timeout until a field is set")
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go srv.Serve(ln)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})
	res, err := http.Get("http://" + ln.Addr().String() + "/files/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || string(body) != "a/b" {
		t.Fatalf("%d %q", res.StatusCode, body)
	}
}

func TestShutdownWaitsThenClosed(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /block", func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: mux}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()

	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		res, err := http.Get("http://" + ln.Addr().String() + "/block")
		if err == nil {
			res.Body.Close()
		}
	}()
	<-started
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("shutdown: %v", err)
	}
	close(release)
	select {
	case <-clientDone:
	case <-time.After(time.Second):
		t.Fatal("client did not return")
	}
	select {
	case err = <-serveErr:
		if err != http.ErrServerClosed {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("Serve did not return")
	}
}

func TestSSEStopsWhenClientCancels(t *testing.T) {
	started := make(chan struct{})
	finished := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "no flush", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = io.WriteString(w, "data: one\n\n")
		flusher.Flush()
		close(started)
		<-r.Context().Done()
		close(finished)
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+ln.Addr().String()+"/events", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	buf := make([]byte, 32)
	n, err := res.Body.Read(buf)
	if n == 0 || err != nil && err != io.EOF {
		t.Fatalf("read %d %v", n, err)
	}
	if string(buf[:n]) != "data: one\n\n" {
		t.Fatalf("event %q", buf[:n])
	}
	cancel()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("handler kept running after the client canceled")
	}
}
