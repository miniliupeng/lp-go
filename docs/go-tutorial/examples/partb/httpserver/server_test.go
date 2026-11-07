package httpserver

import (
	"io"
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
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || string(body) != "ok" {
		t.Fatalf("healthz %d %q", res.StatusCode, body)
	}

	res, err = http.Get(srv.URL + "/items/abc")
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), `"id":"abc"`) && !strings.Contains(string(body), `"id": "abc"`) {
		t.Fatalf("item %d %q", res.StatusCode, body)
	}

	res, err = http.Post(srv.URL+"/echo", "application/json", strings.NewReader(`{"message":"hi"}`))
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), "hi") {
		t.Fatalf("echo %d %q", res.StatusCode, body)
	}

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("POST /healthz = %d", res.StatusCode)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
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
		t.Fatalf("missing id = %q", rec.Header().Get("X-Request-Id"))
	}
}
