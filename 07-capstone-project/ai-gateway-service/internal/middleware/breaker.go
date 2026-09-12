package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateHalfOpen
	StateOpen
)

// CircuitBreaker 自适应断路器中间件 (Closed ⇄ Open ⇄ Half-Open)
type CircuitBreaker struct {
	state             CircuitState
	failureThreshold  int
	openTimeout       time.Duration
	consecutiveFails  int
	lastStateChange   time.Time
	halfOpenSuccesses int
	mu                sync.Mutex
}

func NewCircuitBreaker(failureThreshold int, openTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		openTimeout:      openTimeout,
		lastStateChange:  time.Now(),
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	if cb.state == StateOpen && now.Sub(cb.lastStateChange) > cb.openTimeout {
		cb.state = StateHalfOpen
		cb.lastStateChange = now
		cb.halfOpenSuccesses = 0
	}

	return cb.state != StateOpen
}

func (cb *CircuitBreaker) ReportSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.halfOpenSuccesses++
		if cb.halfOpenSuccesses >= 2 {
			cb.state = StateClosed
			cb.consecutiveFails = 0
			cb.lastStateChange = time.Now()
		}
	} else if cb.state == StateClosed {
		cb.consecutiveFails = 0
	}
}

func (cb *CircuitBreaker) ReportFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.consecutiveFails++
	if cb.state == StateHalfOpen || cb.consecutiveFails >= cb.failureThreshold {
		cb.state = StateOpen
		cb.lastStateChange = time.Now()
	}
}

// statusWriter 用于拦截响应状态码
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (cb *CircuitBreaker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cb.Allow() {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"message": "circuit breaker is open: upstream model degraded",
					"type":    "circuit_breaker_error",
				},
			})
			return
		}

		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)

		if sw.status >= 500 {
			cb.ReportFailure()
		} else {
			cb.ReportSuccess()
		}
	})
}
