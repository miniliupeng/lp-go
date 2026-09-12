package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery 捕获下游任何未处理的 Panic，输出带有堆栈的日志并返回规范的 500 JSON
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				slog.ErrorContext(r.Context(), "server internal panic recovered",
					slog.Any("panic", rec),
					slog.String("stack", stack),
				)

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"message": fmt.Sprintf("internal server error: %v", rec),
						"type":    "internal_error",
					},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
