package main

import (
	"sync"
	"testing"
)

// TestClosedChannelRead 严格单测验证已关闭通道排空后返回零值且 ok=false
func TestClosedChannelRead(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)

	val1, ok1 := <-ch
	if val1 != 42 || !ok1 {
		t.Fatalf("第一次读取应成功获取 42，实际 val=%d, ok=%t", val1, ok1)
	}

	val2, ok2 := <-ch
	if val2 != 0 || ok2 {
		t.Fatalf("第二次读取应返回零值且 ok=false，实际 val=%d, ok=%t", val2, ok2)
	}
}

// TestMutexConcurrentCounter 验证 Mutex 能够彻底避免数据竞争
func TestMutexConcurrentCounter(t *testing.T) {
	var mu sync.Mutex
	var count int
	var wg sync.WaitGroup

	iterations := 500
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()

	if count != iterations {
		t.Fatalf("互斥锁累加预期为 %d，实际为 %d", iterations, count)
	}
}

// BenchmarkUnbufferedChannel 测试无缓冲通道同步交接性能
func BenchmarkUnbufferedChannel(b *testing.B) {
	ch := make(chan int)
	var wg sync.WaitGroup

	b.ReportAllocs()

	go func() {
		for range ch {
		}
	}()

	for b.Loop() {
		ch <- 1
	}
	close(ch)
	wg.Wait()
}

// BenchmarkBufferedChannel 测试有缓冲通道异步收发吞吐
func BenchmarkBufferedChannel(b *testing.B) {
	ch := make(chan int, 100)
	var wg sync.WaitGroup

	b.ReportAllocs()

	go func() {
		for range ch {
		}
	}()

	for b.Loop() {
		ch <- 1
	}
	close(ch)
	wg.Wait()
}
