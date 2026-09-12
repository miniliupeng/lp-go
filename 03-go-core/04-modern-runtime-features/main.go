package main

import (
	"bytes"
	"context"
	"fmt"
	"iter"
	"log/slog"
	"os"
	"runtime/pprof"
)

// ==================== 1. iter.Seq / iter.Seq2 现代迭代器 ====================

// FilterMap 实现通用的惰性流式迭代器 (Zero-Allocation Pipeline)
// 传入基础迭代器，按条件过滤并转换，全程不产生任何临时切片分配
func FilterMap[T any, R any](seq iter.Seq[T], filter func(T) bool, transform func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range seq {
			if filter(v) {
				if !yield(transform(v)) {
					return // 消费方主动 break，短路退出
				}
			}
		}
	}
}

// SliceSeq 将切片包装为标准 iter.Seq
func SliceSeq[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range slice {
			if !yield(item) {
				return
			}
		}
	}
}

// ==================== 2. log/slog 现代高性能结构化日志 ====================

func demoStructuredLogging() {
	var buf bytes.Buffer
	// 创建 JSON 结构化日志 Handler
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 高性能强类型属性注入 (slog.Int, slog.String 避免 interface 装箱)
	logger.InfoContext(context.Background(), "user_login_success",
		slog.Int64("user_id", 10086),
		slog.String("ip", "192.168.1.100"),
		slog.Bool("is_admin", false),
	)

	fmt.Printf("[slog 结构化输出] %s", buf.String())
}

// ==================== 3. PGO (Profile-Guided Optimization) 密集计算载荷 ====================

// 密集业务计算函数：用于模拟线上核心链路并生成 PGO 分析文件
//go:noinline
func simulateHotCompute(x int) int {
	result := 0
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			result += x * 3
		} else {
			result -= x
		}
	}
	return result
}

func demoGeneratePGOProfile() {
	// 生成 default.pgo 样本文件供后续编译器自动化优化
	f, err := os.Create("default.pgo")
	if err != nil {
		fmt.Printf("无法创建 PGO 文件: %v\n", err)
		return
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("启动 CPU Profile 失败: %v\n", err)
		return
	}

	// 运行密集循环，让 profiler 捕获足够的热点调用样本
	total := 0
	for i := 0; i < 50000; i++ {
		total += simulateHotCompute(i)
	}
	pprof.StopCPUProfile()

	fmt.Printf("[PGO 采样完成] 已生成生产级 default.pgo 样本数据，计算校验和: %d\n", total)
}

func main() {
	fmt.Println("=== 1. Go 1.23+ 原生迭代器 (iter.Seq) 惰性流水线 ===")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numSeq := SliceSeq(nums)

	// 构建流水线：筛选偶数 -> 乘以 10 -> 取出并打印
	evenMulti10 := FilterMap(numSeq,
		func(n int) bool { return n%2 == 0 },
		func(n int) int { return n * 10 },
	)

	fmt.Print("流水线输出结果: ")
	for val := range evenMulti10 {
		fmt.Printf("%d ", val)
	}
	fmt.Println()

	fmt.Println("\n=== 2. Go 1.21+ log/slog 结构化日志 ===")
	demoStructuredLogging()

	fmt.Println("\n=== 3. Go 1.20+ PGO 性能剖析画像采集演练 ===")
	demoGeneratePGOProfile()
}
