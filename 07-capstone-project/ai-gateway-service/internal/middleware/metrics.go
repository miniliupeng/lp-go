package middleware

import (
	"net/http"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/pkg/metrics"
)

// Metrics 拦截请求耗时、状态并向 RED 收集器打点
func Metrics(collector *metrics.REDCollector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(sw, r)

			dur := time.Since(start)
			isError := sw.status >= 400
			collector.Record(dur, isError)
		})
	}
}
