package main

import (
	"runtime"
	"testing"
)

// TestEscapeAndGCStats 验证逃逸函数行为与主动 GC 回收能力
func TestEscapeAndGCStats(t *testing.T) {
	p := escapePointer()
	if *p != 42 {
		t.Fatalf("expected 42, got %d", *p)
	}

	closure := escapeClosure()
	if closure() != 1 || closure() != 2 {
		t.Fatal("closure failed to accumulate state")
	}

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 制造可被 GC 的临时大块内存
	temp := make([]byte, 10*1024*1024)
	temp[0] = 1
	_ = temp
	runtime.GC()

	runtime.ReadMemStats(&m2)
	if m2.NumGC <= m1.NumGC {
		t.Fatalf("expected NumGC to increase, before=%d, after=%d", m1.NumGC, m2.NumGC)
	}
}

// 1. 纯栈上分配基准测试 (0 allocs/op, 纳秒级极速)
func BenchmarkStackAllocation(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		res := stackAllocOnly(10, 20)
		_ = res
	}
}

// 2. 堆逃逸分配基准测试 (产生 Heap Alloc 与 GC 追踪开销)
func BenchmarkHeapAllocation(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		p := escapePointer()
		_ = p
	}
}

// 3. 闭包捕获逃逸基准测试
func BenchmarkClosureAllocation(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		fn := escapeClosure()
		_ = fn()
	}
}
