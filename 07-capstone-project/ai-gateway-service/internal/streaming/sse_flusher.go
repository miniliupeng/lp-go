package streaming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
)

var ErrFlusherUnsupported = errors.New("streaming unsupported by underlying connection")

// SSEFlusher 基于 http.Flusher 与 fmt.Appendf 的零多余内存逃逸流式推帧器
type SSEFlusher struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func NewSSEFlusher(w http.ResponseWriter) (*SSEFlusher, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, ErrFlusherUnsupported
	}

	// 生产级 SSE 必要响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // 规避 Nginx 缓冲延迟

	return &SSEFlusher{w: w, flusher: flusher}, nil
}

// PushChunk 将增量 Chunk 格式化为 data: JSON 帧推送到客户端并即刻刷新
func (s *SSEFlusher) PushChunk(chunk *domain.ChatCompletionChunk) error {
	payload, err := json.Marshal(chunk)
	if err != nil {
		return err
	}

	// 使用 fmt.Appendf 避免 []byte(fmt.Sprintf(...)) 的临时 string 堆分配
	frame := fmt.Appendf(nil, "data: %s\n\n", payload)
	if _, err := s.w.Write(frame); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

// PushDone 发送 [DONE] 标识告知客户端推流完毕
func (s *SSEFlusher) PushDone() error {
	if _, err := s.w.Write([]byte("data: [DONE]\n\n")); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

// PushPing 发送注释心跳帧保活
func (s *SSEFlusher) PushPing() error {
	if _, err := s.w.Write([]byte(": ping\n\n")); err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}

// StartHeartbeat 开启后台心跳定时器
func (s *SSEFlusher) StartHeartbeat(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = s.PushPing()
			}
		}
	}()
}
