package main

import (
	"runtime"
	"testing"
	"time"
)

// TestMemoryLeakFixVerification 验证切片泄漏修复机制与协程泄漏检出
func TestMemoryLeakFixVerification(t *testing.T) {
	// 验证深拷贝切片长度与值
	clean := FixedMemoryProducer()
	if len(clean) != 10 || clean[0] != 'F' {
		t.Fatalf("expected 10 bytes with 'F', got %v", clean)
	}

	beforeG := runtime.NumGoroutine()
	LeakGoroutineByBlockedChan()
	time.Sleep(20 * time.Millisecond)
	afterG := runtime.NumGoroutine()
	if afterG <= beforeG {
		t.Fatalf("expected goroutine count to increase after leak, before=%d, after=%d", beforeG, afterG)
	}
}

// 1. 模拟小段 CPU 计算基准
func BenchmarkCPULoad(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		count := 0
		for i := 0; i < 1000; i++ {
			count += i * i
		}
		_ = count
	}
}

// 2. 切片深拷贝防泄漏分配基准
func BenchmarkFixedMemoryProducer(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		res := FixedMemoryProducer()
		_ = res
	}
}
