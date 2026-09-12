package main

import (
	"testing"
)

// TestDeferOrder 严格单测验证 defer 命名与匿名返回值的时序差异
func TestDeferOrder(t *testing.T) {
	if res := deferF1(); res != 5 {
		t.Fatalf("deferF1 预期 5，实际得到 %d", res)
	}
	if res := deferF2(); res != 6 {
		t.Fatalf("deferF2 预期 6，实际得到 %d", res)
	}
	if res := deferF3(); res != 5 {
		t.Fatalf("deferF3 预期 5，实际得到 %d", res)
	}
}

// BenchmarkOpenCodedDefer 测试 Go 1.14+ 开放编码 defer 的调用开销
func BenchmarkOpenCodedDefer(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		func() {
			var x int
			defer func() {
				x++
			}()
		}()
	}
}

// BenchmarkDirectCall 测试无 defer 的纯净函数调用开销
func BenchmarkDirectCall(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		func() {
			var x int
			x++
		}()
	}
}
