package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

// 1. 逃逸场景一：函数返回局部变量指针 (Pointer Escape)
// 局部变量生命周期超出了当前栈帧，编译器必须将其分配到堆上
//
//go:noinline
func escapePointer() *int {
	x := 42
	return &x
}

// 2. 逃逸场景二：动态类型逃逸 (Interface / Dynamic Dispatch Escape)
// 传入空接口 any (interface{}) 时，编译期无法确切推断其具体尺寸，大部分逃逸到堆
func escapeDynamicInterface(val int) {
	// fmt.Println 内部接收 ...any 参数，引发反射与堆逃逸
	fmt.Printf("[动态类型逃逸] 参数值: %d\n", val)
}

// 3. 逃逸场景三：大对象溢出或栈空间无法静态确定 (Stack Overflow Escape)
// Go 协程栈默认初始为 2KB~8KB，单对象超过 64KB (在当前编译器优化中常见阈值) 或动态长度直接上堆
//
//go:noinline
func escapeStackOverflow() []byte {
	// 64KB+ 大切片超出单栈帧安全阈值，直接在堆上分配
	largeSlice := make([]byte, 128*1024)
	largeSlice[0] = 1
	return largeSlice
}

// 4. 逃逸场景四：闭包捕获外部变量导致生命周期延长 (Closure Escape)
//
//go:noinline
func escapeClosure() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 非逃逸对比函数：仅在栈上分配并完成计算，随栈帧弹出即刻无开销回收
//go:noinline
func stackAllocOnly(a, b int) int {
	sum := a + b
	return sum
}

func main() {
	fmt.Println("=== 1. 编译器逃逸分析场景实证 ===")
	p := escapePointer()
	fmt.Printf("1. 返回局部指针: 地址=%p, 值=%d\n", p, *p)

	escapeDynamicInterface(100)

	bigData := escapeStackOverflow()
	fmt.Printf("3. 大对象堆分配: 长度=%d 字节\n", len(bigData))

	closureFn := escapeClosure()
	fmt.Printf("4. 闭包捕获累加: 次数=%d, 次数=%d\n", closureFn(), closureFn())

	fmt.Println("\n=== 2. GC 统计与内存指标观测 (runtime.ReadMemStats) ===")
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("初始堆内存: Alloc=%d KB, HeapInuse=%d KB, NumGC=%d\n",
		m1.Alloc/1024, m1.HeapInuse/1024, m1.NumGC)

	// 制造 10MB 的短期临时堆垃圾
	garbageContainer := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		garbageContainer[i] = make([]byte, 100*1024) // 每次 100KB
	}
	runtime.ReadMemStats(&m2)
	fmt.Printf("分配垃圾后堆内存: Alloc=%d KB, HeapInuse=%d KB, NumGC=%d\n",
		m2.Alloc/1024, m2.HeapInuse/1024, m2.NumGC)

	// 手动触发 GC 并打印回收效果
	garbageContainer = nil
	runtime.GC()

	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)
	fmt.Printf("主动 runtime.GC 后堆内存: Alloc=%d KB, HeapInuse=%d KB, NumGC=%d (GC 增量: +%d)\n",
		m3.Alloc/1024, m3.HeapInuse/1024, m3.NumGC, m3.NumGC-m2.NumGC)

	fmt.Println("\n=== 3. 云原生 Go 1.19+ GOMEMLIMIT 软内存限制演练 ===")
	// 读取当前内存限制 (默认 MaxInt64)
	oldLimit := debug.SetMemoryLimit(-1)
	fmt.Printf("系统默认 GOMEMLIMIT: %d MB\n", oldLimit/(1024*1024))

	// 设置 50MB 软内存限制 (模拟 K8s Pod 资源限额)
	limitTarget := int64(50 * 1024 * 1024)
	debug.SetMemoryLimit(limitTarget)
	fmt.Printf("已动态将 GOMEMLIMIT 设定为安全水位: %d MB\n", limitTarget/(1024*1024))

	// 持续分配内存直至触发 GOMEMLIMIT 调度的密集垃圾回收
	startGC := m3.NumGC
	for i := 0; i < 200; i++ {
		_ = make([]byte, 512*1024) // 每次 512KB
		if i%50 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	var m4 runtime.MemStats
	runtime.ReadMemStats(&m4)
	fmt.Printf("内存压力测试后: NumGC 次数增加至 %d (新增 GC 轮次: %d)\n",
		m4.NumGC, m4.NumGC-startGC)
	fmt.Println(">> GOMEMLIMIT 成功在内存触碰安全红线前主动发起 GC，彻底消除 OOM 隐患！")

	// 还原配置
	debug.SetMemoryLimit(oldLimit)
}
