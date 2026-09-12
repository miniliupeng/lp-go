package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"unsafe"
)

// TestMemoryAlignment 验证内存重排效果与切片深拷贝正确性
func TestMemoryAlignment(t *testing.T) {
	if unsafe.Sizeof(BadStruct{}) != 24 {
		t.Fatalf("expected BadStruct size 24, got %d", unsafe.Sizeof(BadStruct{}))
	}
	if unsafe.Sizeof(GoodStruct{}) != 16 {
		t.Fatalf("expected GoodStruct size 16, got %d", unsafe.Sizeof(GoodStruct{}))
	}

	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	copy(dst, src)
	dst[0] = 999
	if src[0] == 999 {
		t.Fatal("copy should perform deep copy, but src was modified")
	}
}

var testParts = []string{"hello", "world", "golang", "backend", "developer", "high", "concurrency"}

// 1. + 运算符拼接
func BenchmarkConcatPlus(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var s string
		for _, part := range testParts {
			s += part
		}
		_ = s
	}
}

// 2. fmt.Sprintf 格式化拼接
func BenchmarkConcatSprintf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = fmt.Sprintf("%s%s%s%s%s%s%s",
			testParts[0], testParts[1], testParts[2], testParts[3], testParts[4], testParts[5], testParts[6])
	}
}

// 3. strings.Join 拼接
func BenchmarkConcatStringsJoin(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = strings.Join(testParts, "")
	}
}

// 4. bytes.Buffer 缓冲区拼接
func BenchmarkConcatBytesBuffer(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var buf bytes.Buffer
		for _, part := range testParts {
			buf.WriteString(part)
		}
		_ = buf.String()
	}
}

// 5. strings.Builder 零拷贝最优解 (预分配内存)
func BenchmarkConcatStringsBuilder(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var sb strings.Builder
		sb.Grow(64)
		for _, part := range testParts {
			sb.WriteString(part)
		}
		_ = sb.String()
	}
}

// 结构体字段重排基准
func BenchmarkAllocBadStruct(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		slice := make([]BadStruct, 1000)
		_ = slice
	}
}

func BenchmarkAllocGoodStruct(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		slice := make([]GoodStruct, 1000)
		_ = slice
	}
}
