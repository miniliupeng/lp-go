// 并发入门。结果不依赖调度顺序。运行：
//
//	go run ./docs/go-tutorial/examples/parta/concurrency
package main

import (
	"fmt"
	"sync"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}
	results := make(chan int, len(nums))

	var wg sync.WaitGroup
	for _, n := range nums {
		wg.Add(1)
		// 参数 n 再传一次：意图明确，也不依赖某一版循环变量的实现。
		go func(n int) {
			defer wg.Done()
			results <- n * n
		}(n)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for v := range results {
		sum += v
	}
	fmt.Printf("sumOfSquares=%d\n", sum)
}
