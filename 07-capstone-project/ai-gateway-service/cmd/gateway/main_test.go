package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/configs"
	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
	"lp-go/07-capstone-project/ai-gateway-service/internal/middleware"
	"lp-go/07-capstone-project/ai-gateway-service/internal/rag"
	"lp-go/07-capstone-project/ai-gateway-service/internal/service"
	"lp-go/07-capstone-project/ai-gateway-service/internal/streaming"
	"lp-go/07-capstone-project/ai-gateway-service/internal/upstream"
	"lp-go/07-capstone-project/ai-gateway-service/pkg/metrics"
)

func setupTestMux() (*http.ServeMux, *upstream.MockEngine, *metrics.REDCollector) {
	cfg := configs.DefaultConfig()
	collector := metrics.NewREDCollector()
	mockEngine := upstream.NewMockEngine(1 * time.Millisecond)
	chatSvc := service.NewChatService(mockEngine, collector)

	rateLimiter := middleware.NewRateLimiter(10000000, 10000000) // 高阈值确保压测不受 429 干扰
	breaker := middleware.NewCircuitBreaker(cfg.Breaker.FailureThreshold, cfg.Breaker.OpenTimeout)

	mux := http.NewServeMux()

	chatHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req domain.ChatCompletionRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		ctx := r.Context()

		if req.Stream {
			flusher, err := streaming.NewSSEFlusher(w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = chatSvc.Stream(ctx, &req, flusher)
		} else {
			resp, err := chatSvc.Complete(ctx, &req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(resp)
		}
	})

	pipeline := middleware.Recovery(
		middleware.Metrics(collector)(
			middleware.Trace(
				rateLimiter.Middleware(
					breaker.Middleware(chatHandler),
				),
			),
		),
	)

	mux.Handle("/v1/chat/completions", pipeline)
	return mux, mockEngine, collector
}

func TestSyncCompletion(t *testing.T) {
	mux, _, _ := setupTestMux()

	reqBody := domain.ChatCompletionRequest{
		Model: "deepseek-r1-chat",
		Messages: []domain.ChatMessage{
			{Role: domain.RoleUser, Content: "你好，请自我介绍"},
		},
		Stream: false,
	}
	data, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(data))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d", w.Code)
	}

	traceID := w.Header().Get("X-Trace-ID")
	if traceID == "" {
		t.Fatal("expected X-Trace-ID header in response")
	}

	var resp domain.ChatCompletionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(resp.Choices) == 0 || !strings.Contains(resp.Choices[0].Message.Content, "Mock Response") {
		t.Fatalf("unexpected choice content: %+v", resp.Choices)
	}
}

func TestSSEStreamCompletion(t *testing.T) {
	mux, _, _ := setupTestMux()

	reqBody := domain.ChatCompletionRequest{
		Model: "deepseek-r1-chat",
		Messages: []domain.ChatMessage{
			{Role: domain.RoleUser, Content: "流式测试"},
		},
		Stream:    true,
		EnableRAG: true,
	}
	data, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(data))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got: %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("expected text/event-stream, got: %s", w.Header().Get("Content-Type"))
	}

	bodyStr := w.Body.String()
	if !strings.Contains(bodyStr, "data: [DONE]") {
		t.Fatal("expected [DONE] stream marker in SSE output")
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := middleware.NewRateLimiter(1, 1) // 容量 1
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req1 := httptest.NewRequest("POST", "/", nil)
	req1.Header.Set("Authorization", "Bearer key-user-1")
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("first request should pass, got: %d", w1.Code)
	}

	// 紧接着发第二个请求，必然被限流拦截
	req2 := httptest.NewRequest("POST", "/", nil)
	req2.Header.Set("Authorization", "Bearer key-user-1")
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("second burst request should be 429, got: %d", w2.Code)
	}
}

func TestCircuitBreakerTripping(t *testing.T) {
	cb := middleware.NewCircuitBreaker(2, 50*time.Millisecond)
	handler := cb.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Fail") == "true" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	// 连续 2 次 500 触发熔断
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Fail", "true")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	// 此时断路器应直接返回 503 Service Unavailable
	reqTrip := httptest.NewRequest("GET", "/", nil)
	wTrip := httptest.NewRecorder()
	handler.ServeHTTP(wTrip, reqTrip)
	if wTrip.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when circuit breaker open, got: %d", wTrip.Code)
	}
}

func TestRAGHybridSearch(t *testing.T) {
	vecIdx := rag.NewVectorIndex()
	bm25Idx := rag.NewBM25Index()
	merger := rag.NewRRFMerger(60.0)

	doc := &domain.Document{
		ID:        100,
		Title:     "分布式事务",
		Content:   "二阶段提交与 TCC 模式对比",
		Embedding: []float32{0.9, 0.1, 0.0},
		Keywords:  []string{"分布式", "事务", "TCC"},
	}
	vecIdx.Add(doc)
	bm25Idx.Add(doc)

	denseHits := vecIdx.SearchTopK([]float32{0.88, 0.12, 0.0}, 1)
	sparseHits := bm25Idx.SearchTopK([]string{"事务"}, 1)
	finalHits := merger.Merge(denseHits, sparseHits, 1)

	if len(finalHits) == 0 || finalHits[0].Doc.ID != 100 {
		t.Fatalf("expected hit doc ID 100, got: %+v", finalHits)
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	v1 := make([]float32, 1536)
	v2 := make([]float32, 1536)
	for i := 0; i < 1536; i++ {
		v1[i] = float32(i) * 0.001
		v2[i] = float32(1536-i) * 0.001
	}

	b.ResetTimer()
	for b.Loop() {
		_ = rag.CosineSimilarity(v1, v2)
	}
}

func BenchmarkGatewaySyncRequest(b *testing.B) {
	mux, _, _ := setupTestMux()
	reqBody := domain.ChatCompletionRequest{
		Model: "deepseek-r1-chat",
		Messages: []domain.ChatMessage{
			{Role: domain.RoleUser, Content: "基准测试问题"},
		},
		Stream: false,
	}
	data, _ := json.Marshal(reqBody)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(data))
		req.Header.Set("Authorization", "Bearer bench-token")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("expected 200, got %d", w.Code)
		}
	}
}
