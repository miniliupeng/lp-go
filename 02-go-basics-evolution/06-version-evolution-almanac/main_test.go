package main

import (
	"errors"
	"fmt"
	"testing"
)

// TestVersionEvolutionFeatures 验证新版本内置 clear 及现代错误链
func TestVersionEvolutionFeatures(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	clear(m)
	if len(m) != 0 {
		t.Fatalf("expected cleared map len 0, got %d", len(m))
	}

	rawErr := ErrNotFound
	wrapped := fmt.Errorf("wrap: %w", rawErr)
	if !errors.Is(wrapped, ErrNotFound) {
		t.Fatal("errors.Is failed to unwrap target sentinel error")
	}

	// 验证泛型函数
	intSum := SumSlice([]int{1, 2, 3, 4})
	if intSum != 10 {
		t.Fatalf("SumSlice[int] expected 10, got %d", intSum)
	}

	floatSum := SumSlice([]float64{1.1, 2.2})
	if floatSum < 3.29 || floatSum > 3.31 {
		t.Fatalf("SumSlice[float64] expected ~3.3, got %f", floatSum)
	}
}

// BenchmarkBuiltinClearSlice 测试 Go 1.21+ 内置 clear(slice) 的执行速度（底层 memclr 优化）
func BenchmarkBuiltinClearSlice(b *testing.B) {
	b.ReportAllocs()
	s := make([]int, 1024)
	for b.Loop() {
		clear(s)
	}
}

// BenchmarkManualSliceClear 测试传统 for 遍历逐个置零
func BenchmarkManualSliceClear(b *testing.B) {
	b.ReportAllocs()
	s := make([]int, 1024)
	for b.Loop() {
		for i := range s {
			s[i] = 0
		}
	}
}
