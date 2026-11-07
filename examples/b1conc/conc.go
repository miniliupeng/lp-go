// Package conc 是并发入门：固定次数接收，以及 WaitGroup.Go。
package conc

import "sync"

// SumSquares 启动 n 个 goroutine，每个发送一个平方。
// 接收恰好 n 次，所以不需要 close。
func SumSquares(n int) int {
	ch := make(chan int)
	for i := 1; i <= n; i++ {
		go func() { ch <- i * i }()
	}
	sum := 0
	for range n {
		sum += <-ch
	}
	return sum
}

// SumSquaresLocked 用互斥锁累加。Go 1.27 的 WaitGroup.Go 会 Add(1)，函数返回后 Done。
// 文档要求传入的函数不要 panic。Go 1.22 起每次迭代有自己的 i。
func SumSquaresLocked(n int) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sum := 0
	for i := 1; i <= n; i++ {
		wg.Go(func() {
			mu.Lock()
			sum += i * i
			mu.Unlock()
		})
	}
	wg.Wait()
	return sum
}

// First 在两个 channel 里取先就绪的一个。两者都就绪时选择是伪随机的，调用方不要依赖顺序。
func First(data <-chan string, stop <-chan struct{}) (string, bool) {
	select {
	case v := <-data:
		return v, true
	case <-stop:
		return "", false
	}
}
