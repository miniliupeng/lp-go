package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== 1. 新手入门：sync.WaitGroup 优雅并发等待 ===")
	demoWaitGroup()

	fmt.Println("\n=== 2. 基础入门：sync.Mutex 与 sync.RWMutex 并发安全保护 ===")
	demoMutexAndRWMutex()

	fmt.Println("\n=== 3. 硬核底层：Channel 状态矩阵之从 Closed 读取 ===")
	demoClosedChannelRead()

	fmt.Println("\n=== 4. 硬核底层：select 多路复用与 default 非阻塞 ===")
	demoSelectNonBlocking()

	fmt.Println("\n=== 5. 硬核底层：context.Context 树形超时级联取消 ===")
	demoContextCancellation()
}

func demoMutexAndRWMutex() {
	// 1. sync.Mutex 互斥锁保护并发累加
	var mu sync.Mutex
	var count int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("互斥锁保护 1000 次并发累加结果: %d (预期: 1000)\n", count)

	// 2. sync.RWMutex 读写互斥锁（读多写少优化：允许多读，独占单写）
	var rw sync.RWMutex
	cache := make(map[string]string)

	// 写锁独占
	rw.Lock()
	cache["engine"] = "turbofan"
	rw.Unlock()

	// 读锁共享
	rw.RLock()
	val := cache["engine"]
	rw.RUnlock()
	fmt.Printf("读写锁 RLock 安全读取配置: %s\n", val)
}

func demoWaitGroup() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d 启动并执行完毕\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有 Worker 全部平稳完成！")
}

func demoClosedChannelRead() {
	ch := make(chan int, 2)
	ch <- 100
	ch <- 200
	close(ch) // 关闭通道

	// 1. 读取缓冲区内已存在的数据
	val1, ok1 := <-ch
	fmt.Printf("关闭后第 1 次读取: val=%d, ok=%t (读取缓冲内数据)\n", val1, ok1)
	val2, ok2 := <-ch
	fmt.Printf("关闭后第 2 次读取: val=%d, ok=%t (读取缓冲内数据)\n", val2, ok2)

	// 2. 缓冲区排空后继续读取：立即返回零值与 ok=false，绝不阻塞！
	val3, ok3 := <-ch
	fmt.Printf("关闭且排空后第 3 次读取: val=%d (返回类型零值), ok=%t (明确告知已关闭)\n", val3, ok3)
}

func demoSelectNonBlocking() {
	ch := make(chan string, 1)

	// 尝试非阻塞写入
	select {
	case ch <- "ping":
		fmt.Println("成功非阻塞写入 ping")
	default:
		fmt.Println("通道已满，跳入 default 避免阻塞")
	}

	// 再次尝试非阻塞写入（此时已满）
	select {
	case ch <- "pong":
		fmt.Println("成功写入 pong")
	default:
		fmt.Println("通道已满，成功触发 default 非阻塞分支！")
	}
}

func demoContextCancellation() {
	// 创建一个 50ms 后超时的根 Context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动两个子协程监听同一个 Context
	worker := func(id int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d 接收到 Context 取消信号: %v，立即优雅退出！\n", id, ctx.Err())
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	go worker(1)
	go worker(2)

	wg.Wait()
	fmt.Println("所有子协程均被 Context 树形信号平稳回收！")
}
