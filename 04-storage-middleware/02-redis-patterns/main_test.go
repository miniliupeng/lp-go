package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestSingleFlightDedup 验证并发去重与单次回源
func TestSingleFlightDedup(t *testing.T) {
	sf := NewSingleFlightGroup()
	var callCount int64

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := sf.Do("key1", func() (any, error) {
				atomic.AddInt64(&callCount, 1)
				time.Sleep(10 * time.Millisecond)
				return "success_val", nil
			})
			if err != nil || val != "success_val" {
				t.Errorf("unexpected val: %v, err: %v", val, err)
			}
		}()
	}
	wg.Wait()

	if callCount != 1 {
		t.Fatalf("expected single execution (1), but got: %d", callCount)
	}
}

// TestRedisLockWatchdog 验证锁获取与安全释放
func TestRedisLockWatchdog(t *testing.T) {
	lock := NewRedisLock("test_token_999")
	ok := lock.LockWithWatchdog(30 * time.Millisecond)
	if !ok {
		t.Fatal("failed to acquire lock")
	}
	time.Sleep(45 * time.Millisecond) // 触发至少一次续期
	lock.Unlock()
}

// 1. 无 SingleFlight 时的普通并发基准
func BenchmarkWithoutSingleFlight(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		// 模拟直接访问后端计算
		data := "cached_payload"
		_ = data
	}
}

// 2. SingleFlight 请求合并高并发基准测试
func BenchmarkSingleFlightDo(b *testing.B) {
	sf := NewSingleFlightGroup()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			val, _ := sf.Do("bench_key", func() (any, error) {
				return "bench_val", nil
			})
			_ = val
		}
	})
}
