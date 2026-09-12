package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"
)

// TestIteratorAndSlog 验证现代迭代器流水线与 slog 结构化能力
func TestIteratorAndSlog(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	seq := SliceSeq(nums)
	var collected []int

	for v := range FilterMap(seq, func(x int) bool { return x%2 == 0 }, func(x int) int { return x * 10 }) {
		collected = append(collected, v)
	}

	if len(collected) != 3 || collected[0] != 20 || collected[1] != 40 || collected[2] != 60 {
		t.Fatalf("unexpected pipeline result: %v", collected)
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Info("ping", slog.String("key", "val"))
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("expected json key val, got: %s", buf.String())
	}
}

// 1. 传统物化切片方式 (每次需要分配新切片存放过滤结果)
func BenchmarkTraditionalSliceFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ReportAllocs()
	for b.Loop() {
		var res []int
		for _, v := range data {
			if v%2 == 0 {
				res = append(res, v*10)
			}
		}
		_ = res
	}
}

// 2. 原生 iter.Seq 零分配惰性流水线 (0 allocs/op)
func BenchmarkIteratorPipeline(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ReportAllocs()
	for b.Loop() {
		seq := SliceSeq(data)
		pipeline := FilterMap(seq, func(v int) bool { return v%2 == 0 }, func(v int) int { return v * 10 })
		sum := 0
		for v := range pipeline {
			sum += v
		}
		_ = sum
	}
}

// 3. log/slog 高性能结构化属性注入吞吐测试
func BenchmarkSlogTextHandler(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	ctx := context.Background()

	b.ReportAllocs()
	for b.Loop() {
		logger.InfoContext(ctx, "bench_order_event",
			slog.Int64("order_id", 888888),
			slog.String("status", "PAID"),
			slog.Float64("amount", 99.5),
		)
	}
}
