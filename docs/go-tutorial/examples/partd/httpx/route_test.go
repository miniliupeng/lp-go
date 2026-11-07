package httpx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestZeroTimeoutsMeanNoDeadline(t *testing.T) {
	// 文档：这些字段为 0 或负数时表示不设对应的超时。
	// ReadHeaderTimeout 为 0 时会改用 ReadTimeout；两者都是 0 才是没有超时。
	s := &http.Server{}
	if s.ReadTimeout != 0 || s.ReadHeaderTimeout != 0 || s.WriteTimeout != 0 || s.IdleTimeout != 0 {
		t.Fatalf("%v %v %v %v", s.ReadTimeout, s.ReadHeaderTimeout, s.WriteTimeout, s.IdleTimeout)
	}
	if s.ReadHeaderTimeout != time.Duration(0) {
		t.Fatal(s.ReadHeaderTimeout)
	}
}

func TestWildcardPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, r.PathValue("path"))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/files/a/b")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || string(body) != "a/b" {
		t.Fatalf("%d %q", res.StatusCode, body)
	}
}
