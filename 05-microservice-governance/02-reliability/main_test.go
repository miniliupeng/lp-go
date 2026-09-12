package main

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenBucketLimiter(t *testing.T) {
	limiter := NewTokenBucketLimiter(100, 10)

	// 前 10 个请求必须直接通过
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Fatalf("request %d should be allowed", i)
		}
	}

	// 第 11 个应立即被拒绝（桶已空且未到补充间隔）
	if limiter.Allow() {
		t.Fatal("request 11 should be rejected")
	}

	// 等待 20ms，100/s 应该补充 2 个令牌
	time.Sleep(25 * time.Millisecond)
	if !limiter.Allow() {
		t.Fatal("request after 25ms should be allowed")
	}
}

func TestCircuitBreakerStateTransitions(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected closed, got %v", cb.CurrentState())
	}

	// 触发 2 次错误
	_ = cb.Execute(func() error { return errors.New("err1") })
	_ = cb.Execute(func() error { return errors.New("err2") })

	if cb.CurrentState() != StateOpen {
		t.Fatalf("expected open, got %v", cb.CurrentState())
	}

	// 立即调用应该拦截
	err := cb.Execute(func() error { return nil })
	if !errors.Is(err, ErrCircuitBreakerOpen) {
		t.Fatalf("expected ErrCircuitBreakerOpen, got %v", err)
	}

	// 超时恢复
	time.Sleep(60 * time.Millisecond)
	// 连续 2 次成功
	_ = cb.Execute(func() error { return nil })
	_ = cb.Execute(func() error { return nil })

	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected closed after recovery, got %v", cb.CurrentState())
	}
}

func TestIdempotency(t *testing.T) {
	im := NewIdempotencyManager(time.Minute)
	var callCount atomic.Int32

	action := func() (string, error) {
		callCount.Add(1)
		return "RESULT_OK", nil
	}

	res1, err1 := im.ExecuteIdempotent("KEY_123", action)
	if err1 != nil || res1 != "RESULT_OK" {
		t.Fatalf("res1 error: %v", err1)
	}

	res2, err2 := im.ExecuteIdempotent("KEY_123", action)
	if err2 != nil || res2 != "RESULT_OK" {
		t.Fatalf("res2 error: %v", err2)
	}

	if callCount.Load() != 1 {
		t.Fatalf("action should be called exactly once, got %d", callCount.Load())
	}
}

func BenchmarkTokenBucketLimiter(b *testing.B) {
	limiter := NewTokenBucketLimiter(10000000, 10000000)
	b.ResetTimer()
	for b.Loop() {
		_ = limiter.Allow()
	}
}

func BenchmarkIdempotencyManager(b *testing.B) {
	im := NewIdempotencyManager(time.Minute)
	_, _ = im.ExecuteIdempotent("BENCH_KEY", func() (string, error) {
		return "OK", nil
	})

	b.ResetTimer()
	for b.Loop() {
		_, _ = im.ExecuteIdempotent("BENCH_KEY", func() (string, error) {
			return "OK", nil
		})
	}
}
