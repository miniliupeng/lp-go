package middleware

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type tokenBucket struct {
	rate       float64
	capacity   float64
	tokens     float64
	lastUpdate time.Time
	mu         sync.Mutex
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.lastUpdate = now
	tb.tokens = math.Min(tb.capacity, tb.tokens+elapsed*tb.rate)

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

// RateLimiter 基于租户维度 (API-Key 或 IP) 的惰性令牌桶限流中间件
type RateLimiter struct {
	rate     float64
	capacity float64
	buckets  sync.Map
}

func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		capacity: capacity,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 优先从 Authorization 提取 API-Key，无则按 RemoteAddr
		key := r.Header.Get("Authorization")
		if key == "" {
			key = r.RemoteAddr
			if idx := strings.LastIndex(key, ":"); idx != -1 {
				key = key[:idx]
			}
		}

		val, _ := rl.buckets.LoadOrStore(key, &tokenBucket{
			rate:       rl.rate,
			capacity:   rl.capacity,
			tokens:     rl.capacity,
			lastUpdate: time.Now(),
		})
		bucket := val.(*tokenBucket)

		if !bucket.allow() {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"message": "rate limit exceeded: quota depleted",
					"type":    "rate_limit_error",
				},
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
