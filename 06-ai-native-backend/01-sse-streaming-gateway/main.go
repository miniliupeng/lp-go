package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 1. SSE (Server-Sent Events) 事件协议定义
// ============================================================================

// StreamEvent 表示大模型流式输出的一个片段
type StreamEvent struct {
	ID        int    `json:"id"`
	DeltaText string `json:"delta_text"`
	Finished  bool   `json:"finished"`
}

// FormatSSE 将事件格式化为标准 text/event-stream 协议帧
func FormatSSE(event *StreamEvent) []byte {
	payload, _ := json.Marshal(event)
	return fmt.Appendf(nil, "id: %d\nevent: message\ndata: %s\n\n", event.ID, payload)
}

// ============================================================================
// 2. 带反压（Backpressure）与保活心跳的流式通道
// ============================================================================

// BoundedStreamChannel 具备缓冲上限与阻塞反压的流式管道
type BoundedStreamChannel struct {
	ch       chan *StreamEvent
	capacity int
	closed   atomic.Bool
	mu       sync.Mutex
}

func NewBoundedStreamChannel(capacity int) *BoundedStreamChannel {
	return &BoundedStreamChannel{
		ch:       make(chan *StreamEvent, capacity),
		capacity: capacity,
	}
}

// Push 尝试向管道推送事件，若下游消费过慢且超时则触发反压错误
func (b *BoundedStreamChannel) Push(ctx context.Context, event *StreamEvent, timeout time.Duration) error {
	if b.closed.Load() {
		return errors.New("stream channel is closed")
	}

	pushCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-pushCtx.Done():
		return fmt.Errorf("backpressure triggered: downstream is too slow: %w", pushCtx.Err())
	case b.ch <- event:
		return nil
	}
}

func (b *BoundedStreamChannel) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.closed.Swap(true) {
		close(b.ch)
	}
}

func (b *BoundedStreamChannel) Reader() <-chan *StreamEvent {
	return b.ch
}

// ============================================================================
// 3. 生产级 SSE 网关 Handler
// ============================================================================

type SSEGatewayHandler struct {
	channelCap  int
	pushTimeout time.Duration
}

func NewSSEGatewayHandler(channelCap int, pushTimeout time.Duration) *SSEGatewayHandler {
	return &SSEGatewayHandler{
		channelCap:  channelCap,
		pushTimeout: pushTimeout,
	}
}

// HandleChatStream 模拟对接后端 LLM 并以 SSE 协议流式推送给客户端
func (h *SSEGatewayHandler) HandleChatStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// 必须的 SSE 协议头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // 禁用 Nginx 缓冲

	ctx := r.Context()
	streamChan := NewBoundedStreamChannel(h.channelCap)

	// 启动后台协程模拟从上游 LLM 生成 Token
	go func() {
		defer streamChan.Close()
		tokens := []string{"你", "好", "！", "我", "是", "高", "性", "能", "Go", "AI", "网", "关", "。"}
		for i, tok := range tokens {
			select {
			case <-ctx.Done():
				return // 客户端主动断开，终止上游计算
			default:
				evt := &StreamEvent{
					ID:        i + 1,
					DeltaText: tok,
					Finished:  i == len(tokens)-1,
				}
				if err := streamChan.Push(ctx, evt, h.pushTimeout); err != nil {
					return
				}
				time.Sleep(10 * time.Millisecond) // 模拟推理生成间隔
			}
		}
	}()

	// 保活心跳定时器 (防止移动端代理中途掐断连接)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// 客户端断开连接
			return
		case <-ticker.C:
			// 发送 SSE 注释帧保活心跳
			_, _ = w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		case evt, ok := <-streamChan.Reader():
			if !ok {
				return
			}
			_, _ = w.Write(FormatSSE(evt))
			flusher.Flush() // 核心：即刻刷新网络缓冲区将 Token 触达前端
			if evt.Finished {
				return
			}
		}
	}
}

func main() {
	fmt.Println("=== 阶段四 专题01：AI 原生后端（SSE 流式网关与反压机制） ===")

	handler := NewSSEGatewayHandler(5, 50*time.Millisecond)

	// 模拟流式推送管道与反压
	streamChan := NewBoundedStreamChannel(2) // 极小容量模拟反压

	// 生产者连续塞入数据
	ctx := context.Background()
	_ = streamChan.Push(ctx, &StreamEvent{ID: 1, DeltaText: "Token 1"}, 10*time.Millisecond)
	_ = streamChan.Push(ctx, &StreamEvent{ID: 2, DeltaText: "Token 2"}, 10*time.Millisecond)

	// 此时管道已满，第三个写入将超时触发反压
	err := streamChan.Push(ctx, &StreamEvent{ID: 3, DeltaText: "Token 3"}, 20*time.Millisecond)
	fmt.Printf("[反压实测] 队列满且下游不消费时的行为: %v\n", err)

	// 释放资源
	streamChan.Close()
	_ = handler
}
