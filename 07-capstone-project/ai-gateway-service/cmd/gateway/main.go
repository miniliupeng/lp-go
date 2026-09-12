package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/configs"
	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
	"lp-go/07-capstone-project/ai-gateway-service/internal/middleware"
	"lp-go/07-capstone-project/ai-gateway-service/internal/service"
	"lp-go/07-capstone-project/ai-gateway-service/internal/streaming"
	"lp-go/07-capstone-project/ai-gateway-service/internal/upstream"
	"lp-go/07-capstone-project/ai-gateway-service/pkg/logger"
	"lp-go/07-capstone-project/ai-gateway-service/pkg/metrics"
)

func main() {
	// 1. 初始化结构化日志
	log := logger.InitLogger()
	log.Info("initializing enterprise ai-gateway-service...")

	cfg := configs.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		log.Error("invalid configuration", slog.Any("error", err))
		os.Exit(1)
	}

	// 2. 初始化核心组件
	collector := metrics.DefaultCollector
	mockEngine := upstream.NewMockEngine(cfg.Upstream.MockDelay)
	chatSvc := service.NewChatService(mockEngine, collector)

	// 3. 构建洋葱中间件
	rateLimiter := middleware.NewRateLimiter(cfg.Limiter.Rate, cfg.Limiter.Capacity)
	breaker := middleware.NewCircuitBreaker(cfg.Breaker.FailureThreshold, cfg.Breaker.OpenTimeout)

	mux := http.NewServeMux()

	// 核心业务处理 Handler (兼容 OpenAI /v1/chat/completions)
	chatHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		if req.Model == "" {
			req.Model = cfg.Upstream.DefaultModel
		}

		ctx := r.Context()

		if req.Stream {
			// 流式 SSE 分支
			flusher, err := streaming.NewSSEFlusher(w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			flusher.StartHeartbeat(ctx, 15*time.Second)

			if err := chatSvc.Stream(ctx, &req, flusher); err != nil {
				slog.ErrorContext(ctx, "stream transfer failed", slog.Any("error", err))
			}
		} else {
			// 同步非流式分支
			resp, err := chatSvc.Complete(ctx, &req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(resp)
		}
	})

	// 指标查询 Endpoint
	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total_requests": collector.TotalRequests.Load(),
			"total_errors":   collector.TotalErrors.Load(),
			"error_rate_pct": collector.ErrorRate(),
			"p99_duration_ms": collector.P99Duration(),
			"avg_ttft_ms":    collector.AvgTTFT(),
		})
	})

	// 注册路由并通过洋葱中间件装配 (Recovery -> Metrics -> Trace -> RateLimit -> CircuitBreaker -> Handler)
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
	mux.Handle("/metrics", metricsHandler)

	// 4. 构造 HTTP 服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 5. 启动后台服务协程
	go func() {
		log.Info("gateway server started successfully", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server fatal error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// 6. 监听系统中断信号实现生产级优雅关机 (Graceful Shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gateway server gracefully...")

	drainCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.DrainTimeout)
	defer cancel()

	if err := server.Shutdown(drainCtx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	} else {
		log.Info("gateway server exited cleanly, all active connections drained")
	}
}
