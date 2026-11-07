# 11 并发入门

这一章只建立能写对的最小模型。channel 的关闭规则、select、内存模型和锁的细节在 [Part C](../part-c/04-channels.md)。

## 这一章你要会什么

启动 goroutine，用 channel 把结果送回来，用 `WaitGroup` 等它们结束，并且不数据竞争。

## goroutine

`go f()` 会在新的 goroutine 里调用 `f`。`go` 语句本身 happens before 这次调用开始。goroutine 结束**不会**自动 happens before 其他 goroutine 里的任何操作。所以「启动一个 goroutine 改全局变量，main 里马上读」是数据竞争。要观察到结果，用 channel、锁，或 `WaitGroup`。

`main` 返回时进程退出，还在跑的 goroutine 会被直接停掉。

## channel

`ch := make(chan int)` 是无缓冲 channel：发送和接收必须碰面。`make(chan int, n)` 有 n 个缓冲。

发送 `ch <- v`，接收 `v := <-ch` 或 `v, ok := <-ch`。

示例把平方送进一个有缓冲的 channel，缓冲大小等于任务数，因此发送不会被卡死。再单独起一个 goroutine：等 `WaitGroup` 结束后 `close(results)`。`for v := range results` 会在 channel 被关闭且取空之后停下来。

关闭是发送方的责任。这个例子里，唯一知道「不会再发送」的是那个等待 `WaitGroup` 的 goroutine。

## WaitGroup

```go
wg.Add(1)
go func(n int) {
    defer wg.Done()
    results <- n * n
}(n)
wg.Wait()
```

`Add` 要在 `Wait` 之前发生，并且要在启动对应 goroutine 之前或由已经计入的 goroutine 调用。把 `Add` 放进子 goroutine 的开头，可能 `Wait` 先返回。

Go 1.27 的 `sync.WaitGroup` 还有 `Go(func())` 方法：它内部 `Add(1)` 并在函数返回后 `Done`。传入的函数不能 panic。入门示例把 `Add` / `Done` 写开，是为了看见计数。

## 循环变量

示例把 `n` 当作参数传入。Go 1.22 起每次迭代已经有独立的循环变量，但参数写法在任何受支持的版本上都是明确的。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/concurrency
```

输出固定为：

```text
sumOfSquares=55
```

1²+…+5² = 55。channel 的接收顺序不稳定，求和与顺序无关，所以结果稳定。

## 数据竞争

两个 goroutine 同时写同一个普通变量，就是数据竞争。检查：

```bash
go test -race ./docs/go-tutorial/examples/partc/syncx
```

`-race` 会让程序变慢、更吃内存，但能在测试里抓到这类错误。有竞争的程序，语言不保证你读到的是「最新的那个值」。

## 练习

改成计算 `1+2+…+5`，仍然每个数一个 goroutine。先在纸上写出谁负责 `close`。
