package httpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	srv := httptest.NewServer(NewMux())
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	body := read(t, res)
	if res.StatusCode != http.StatusOK || body != "ok" {
		t.Fatalf("%d %q", res.StatusCode, body)
	}

	res, err = http.Post(srv.URL+"/healthz", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("POST /healthz = %d", res.StatusCode)
	}

	res, err = http.Get(srv.URL + "/items/abc")
	if err != nil {
		t.Fatal(err)
	}
	body = read(t, res)
	if res.StatusCode != http.StatusOK || !strings.Contains(body, `"id":"abc"`) {
		t.Fatalf("%d %q", res.StatusCode, body)
	}
}

func TestRequestID(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-Id", "req-9")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent || rec.Header().Get("X-Request-Id") != "req-9" {
		t.Fatalf("%d %q", rec.Code, rec.Header().Get("X-Request-Id"))
	}
	req = httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Header().Get("X-Request-Id") != "none" {
		t.Fatal(rec.Header().Get("X-Request-Id"))
	}
}

func read(t *testing.T, res *http.Response) string {
	t.Helper()
	defer res.Body.Close()
	buf := make([]byte, 256)
	n, _ := res.Body.Read(buf)
	return string(buf[:n])
}
