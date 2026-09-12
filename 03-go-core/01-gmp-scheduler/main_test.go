package main

import (
	"runtime"
	"sync"
	"testing"
)

// TestGMPEnvironment 验证逻辑 CPU 数与 Gosched 正常执行
func TestGMPEnvironment(t *testing.T) {
	if runtime.NumCPU() <= 0 {
		t.Fatal("expected positive CPU count")
	}
	runtime.Gosched()
}

// BenchmarkGoroutineCreation 测试高并发下 Goroutine 创建与上下文开销
func BenchmarkGoroutineCreation(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
		}()
		wg.Wait()
	}
}

// BenchmarkGoschedOverhead 测试主动出让 CPU (runtime.Gosched) 的调度损耗
func BenchmarkGoschedOverhead(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		runtime.Gosched()
	}
}

// BenchmarkConcurrencyThroughput 测试多个 Goroutine 在 P 之间的并发吞吐
func BenchmarkConcurrencyThroughput(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var sum int
		for pb.Next() {
			sum++
		}
	})
}
