package streaming

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
)

var ErrBackpressureTriggered = errors.New("backpressure triggered: consumer is too slow")

// BackpressureChannel 有界缓冲流式管道，提供超时丢弃与断流保护
type BackpressureChannel struct {
	ch       chan *domain.ChatCompletionChunk
	capacity int
	closed   atomic.Bool
	mu       sync.Mutex
}

func NewBackpressureChannel(capacity int) *BackpressureChannel {
	if capacity <= 0 {
		capacity = 16
	}
	return &BackpressureChannel{
		ch:       make(chan *domain.ChatCompletionChunk, capacity),
		capacity: capacity,
	}
}

// Push 尝试向管道推送 Chunk，当下游消费者拥塞且等待超过 timeout 时返回反压错误
func (bc *BackpressureChannel) Push(ctx context.Context, chunk *domain.ChatCompletionChunk, timeout time.Duration) error {
	if bc.closed.Load() {
		return errors.New("channel is already closed")
	}

	pushCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-pushCtx.Done():
		return fmt.Errorf("%w: %v", ErrBackpressureTriggered, pushCtx.Err())
	case bc.ch <- chunk:
		return nil
	}
}

func (bc *BackpressureChannel) Close() {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if !bc.closed.Swap(true) {
		close(bc.ch)
	}
}

func (bc *BackpressureChannel) Out() <-chan *domain.ChatCompletionChunk {
	return bc.ch
}
