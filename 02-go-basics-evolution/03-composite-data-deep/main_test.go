package main

import (
	"testing"
	"unsafe"
)

// TestFullSliceExpression 验证三索引切片彻底防止原数组被篡改
func TestFullSliceExpression(t *testing.T) {
	origin := []int{1, 2, 3, 4, 5}
	sub := origin[1:3:3]
	sub = append(sub, 999)

	// 确认原数组第 4 个元素依然是 4
	if origin[3] != 4 {
		t.Fatalf("原数组元素被意外篡改！预期 4，实际为 %d", origin[3])
	}
}

// BenchmarkTraditionalStringToBytes 测试传统内置强制转换性能
func BenchmarkTraditionalStringToBytes(b *testing.B) {
	text := "这是用于测试高吞吐字符串转换的基准长文本数据内容示范"
	b.ReportAllocs()

	for b.Loop() {
		bs := []byte(text)
		_ = bs
	}
}

// BenchmarkZeroCopyStringToBytes 测试 Go 1.20+ 官方零拷贝性能
func BenchmarkZeroCopyStringToBytes(b *testing.B) {
	text := "这是用于测试高吞吐字符串转换的基准长文本数据内容示范"
	b.ReportAllocs()

	for b.Loop() {
		bs := unsafe.Slice(unsafe.StringData(text), len(text))
		_ = bs
	}
}
