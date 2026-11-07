// Package httpserver 是一个可以用 httptest 测完的最小 HTTP API。
package httpserver

import (
	"encoding/json"
	"io"
	"net/http"
)

// Item 是 GET /items/{id} 的响应。
type Item struct {
	ID string `json:"id"`
}

// Echo 是 POST /echo 的请求和响应。
type Echo struct {
	Message string `json:"message"`
}

// NewMux 注册健康检查、按 id 取资源、以及回显。
// 路由写法要求 Go 1.22 及以上的 ServeMux。
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	})
	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, Item{ID: r.PathValue("id")})
	})
	mux.HandleFunc("POST /echo", func(w http.ResponseWriter, r *http.Request) {
		var req Echo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, req)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
