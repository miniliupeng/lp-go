package middleware

import (
	"net/http"

	"lp-go/07-capstone-project/ai-gateway-service/pkg/trace"
)

// Trace 拦截请求，解析或派生 W3C traceparent 并注入请求 Context 及响应头
func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.ExtractHTTP(r.Header)
		if span == nil {
			span = trace.NewRootSpan()
		} else {
			span = span.ChildSpan()
		}

		ctx := trace.WithContext(r.Context(), span)
		span.InjectHTTP(w.Header())
		w.Header().Set("X-Trace-ID", span.TraceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
