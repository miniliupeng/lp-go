package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// ==================== 1. 致命场景一：CPU 100% 密集计算 ====================

// SimulateCPUSpike 模拟线上死循环或极慢正则/高频哈希运算导致的 CPU 飙高
//go:noinline
func SimulateCPUSpike(duration time.Duration) {
	deadline := time.Now().Add(duration)
	count := 0
	for time.Now().Before(deadline) {
		// 密集无效计算
		count += 1
		_ = count * count
	}
}

// ==================== 2. 致命场景二：切片截取导致的内存泄漏 ====================

var leakedGlobalSlice [][]byte

// LeakingMemoryProducer 截取大数组的前几个字节，导致大数组永久无法被 GC 回收
//go:noinline
func LeakingMemoryProducer() []byte {
	// 开辟 10MB 的底层大数组
	largeArray := make([]byte, 10*1024*1024)
	largeArray[0] = 'L'
	// 仅截取前 10 个字节，但其底层 Data 指针仍指向整块 10MB 内存！
	return largeArray[:10]
}

// FixedMemoryProducer 正确做法：使用 copy 将所需数据深拷贝出来，使大数组可被 GC
//go:noinline
func FixedMemoryProducer() []byte {
	largeArray := make([]byte, 10*1024*1024)
	largeArray[0] = 'F'
	cleanCopy := make([]byte, 10)
	copy(cleanCopy, largeArray[:10])
	return cleanCopy
}

// ==================== 3. 致命场景三：Goroutine 悬挂泄漏 ====================

// LeakGoroutineByBlockedChan 模拟发送端因无接收方而永久阻塞导致的协程泄漏
func LeakGoroutineByBlockedChan() {
	ch := make(chan int) // 无缓冲 Channel
	go func() {
		// 永远没有接收者，协程栈和上下文永久悬挂在内存中
		ch <- 1
	}()
}

func main() {
	fmt.Println("=== 1. 采集 CPU Profile 现场样本 ===")
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		fmt.Printf("创建 cpu.pprof 失败: %v\n", err)
		return
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err == nil {
		SimulateCPUSpike(200 * time.Millisecond)
		pprof.StopCPUProfile()
		fmt.Println("[CPU Profile] 样本采集成功 -> 已生成 cpu.pprof")
	}

	fmt.Println("\n=== 2. 模拟内存泄漏场景并采集 Heap Profile ===")
	// 模拟 3 次泄漏产生 30MB 悬挂内存
	for i := 0; i < 3; i++ {
		leakedGlobalSlice = append(leakedGlobalSlice, LeakingMemoryProducer())
	}

	memFile, err := os.Create("mem.pprof")
	if err != nil {
		fmt.Printf("创建 mem.pprof 失败: %v\n", err)
		return
	}
	defer memFile.Close()

	runtime.GC() // 主动触发一次 GC，此时被引用的 30MB 大数组依然无法被释放
	if err := pprof.WriteHeapProfile(memFile); err == nil {
		fmt.Println("[Heap Profile] 样本采集成功 -> 已生成 mem.pprof (包含 30MB 泄漏切片)")
	}

	fmt.Println("\n=== 3. 模拟 Goroutine 悬挂并采集 Goroutine Profile ===")
	initialG := runtime.NumGoroutine()
	for i := 0; i < 20; i++ {
		LeakGoroutineByBlockedChan()
	}
	time.Sleep(50 * time.Millisecond)
	currentG := runtime.NumGoroutine()
	fmt.Printf("[Goroutine 泄漏观测] 初始协程数: %d -> 当前泄漏后协程数: %d (泄漏增量: +%d)\n",
		initialG, currentG, currentG-initialG)

	gFile, err := os.Create("goroutine.pprof")
	if err != nil {
		fmt.Printf("创建 goroutine.pprof 失败: %v\n", err)
		return
	}
	defer gFile.Close()

	if p := pprof.Lookup("goroutine"); p != nil {
		_ = p.WriteTo(gFile, 0)
		fmt.Println("[Goroutine Profile] 样本采集成功 -> 已生成 goroutine.pprof")
	}
}
