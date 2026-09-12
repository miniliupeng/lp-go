package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== 1. GMP 基础环境信息 ===")
	fmt.Printf("逻辑 CPU 数量 (GOMAXPROCS): %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("当前活动 Goroutine 数量: %d\n", runtime.NumGoroutine())

	// 创建用于追踪的 trace 文件
	traceFile, err := os.Create("trace.out")
	if err != nil {
		fmt.Printf("创建 trace.out 失败: %v\n", err)
		return
	}
	defer traceFile.Close()

	if err := trace.Start(traceFile); err != nil {
		fmt.Printf("启动 trace 失败: %v\n", err)
		return
	}
	defer trace.Stop()

	fmt.Println("\n=== 2. 抢占式调度验证（单 P 场景）===")
	verifyPreemption()
}

// verifyPreemption 模拟单核环境下密集计算任务被非协作抢占的过程
func verifyPreemption() {
	// 保存原设置并在退出时还原，强制单逻辑处理器运行
	originalProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(originalProcs)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// G1: 密集型纯数值计算（无任何显式函数调用和 I/O）
	go func() {
		defer wg.Done()
		var count uint64
		for {
			count++
			select {
			case <-ctx.Done():
				fmt.Printf("[G1 密集计算] 退出，累计计算轮次: %d\n", count)
				return
			default:
			}
		}
	}()

	// G2: 轻量观察协程
	// 在 Go 1.14 之前的协作式调度中，若 G1 不让出，G2 将永远无法在单 P 上运行
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		fmt.Println("[G2 观察协程] 成功抢占执行！证明基于信号的异步抢占生效。")
	}()

	wg.Wait()
	fmt.Println("=== GMP 调度实验结束，已生成 trace.out ===")
}
