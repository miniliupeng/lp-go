# 05 sync 与原子操作

源码依据：`sync` 包文档、`internal/sync/mutex.go` 的饥饿模式注释、`sync/atomic` 包文档、`sync.Map` 文档。示例在 Go 1.27.1 上测试。

## Mutex

零值是未加锁的 Mutex。第一次使用之后不能复制。

`internal/sync/mutex.go` 描述两种模式：

- **正常模式**：等待者按 FIFO 排队，但被唤醒的等待者并不直接拥有锁，它要和新来的 goroutine 竞争。新来的已经在 CPU 上，可能连续抢到。抢失败的等待者放到队首。
- **饥饿模式**：某个等待者超过 1ms（`starvationThresholdNs = 1e6`）还拿不到锁，就进入饥饿模式。解锁时把锁直接交给队首等待者。若一个被唤醒的等待者发现自己是最后一个，或者等待时间已经短于阈值，就退出饥饿模式。

这把锁不可重入。同一 goroutine 再 `Lock` 一次会把自己堵死。

`TryLock` 的文档说：正确用法存在，但很少见；想用 `TryLock` 往往说明锁的用法有问题。失败的 `TryLock` 不建立任何 synchronized-before。

`RWMutex` 允许多个读者或一个写者。读锁不能升级成写锁，升级会死锁。内存模型对 `RLock` / `RUnlock` 和写锁的配对有单独的条款，见 [内存一章](01-memory-escape.md)。

## WaitGroup、Once、Pool

`Add` 必须在对应的 `Wait` 观察到之前完成。把 `Add` 放进子 goroutine 里，可能 `Wait` 先跑完。计数不能变成负数，否则 panic。复用 `WaitGroup` 时，下一轮 `Add` 必须发生在上一轮 `Wait` 已经返回之后。

Go 1.27 提供 `WaitGroup.Go(f)`：内部 `Add(1)`，`f` 返回后 `Done`。文档要求 `f` 不能 panic。示例 [sync_test.go](../examples/partc/syncx/sync_test.go) 用它做计数。

`Once.Do` 保证 `f` 只执行一次。并发的 `Do` 会等到这次 `f` 返回。`f` 里再调用同一个 `Once` 会死锁。

`Pool` 用来复用临时对象，减轻 GC。`New` 只在 `Get` 取空时调用。池里的对象在 GC 时可以被丢掉，不能当作稳定缓存。`Put` 进去的指针不要再从外面用。

## 原子操作

`sync/atomic` 的文档：若 B 观察到 A 的效果，则 A synchronized before B；所有原子操作表现得有一个顺序一致的总顺序。这和 C++ 的 SC 原子、Java 的 `volatile` 相同。

优先用 `atomic.Int64` 这种类型，而不是对一个裸 `int64` 调 `atomic.AddInt64`。在 32 位平台上，裸的 64 位原子要求调用方自己保证 64 位对齐；`atomic.Int64` 会处理对齐。

原子加法适合计数器。多个字段要一起变更时，原子变量不够，用 Mutex 保护整组字段。

## sync.Map

文档写得很明确：大多数代码应该用普通 map 加锁。`sync.Map` 针对两种情况做了优化：（1）某个键写一次、读很多次，例如只增长的缓存；（2）不同 goroutine 读写的键集合互不相交。它不是「并发 map 的更快替代」。零值可用，不能复制。一次写 synchronized before 观察到这次写的读。

## 运行

```bash
go test ./docs/go-tutorial/examples/partc/syncx -count=1
go test -race ./docs/go-tutorial/examples/partc/syncx -count=1
```

`Pool` 的测试只断言空池会调用 `New`。它不断言 `Put` 之后 `Get` 一定拿回同一个对象。

## 面试常问

- 正常模式和饥饿模式差在哪，阈值是多少：1ms，直接交接。
- `WaitGroup` 为什么 `Add` 不能放进 goroutine 开头。
- `sync.Map` 适合什么，不适合什么：按文档的两种访问模式回答，不要说它普遍更快。
- 原子操作的内存序：观察到效果就 synchronized before，并且整体顺序一致。不要编一套 Go 里不存在的 Acquire/Release API 名称；标准库的原子就是这一套 SC 语义。

## 对比陷阱

- 用 `Pool` 做连接池或会话缓存。对象会消失。
- 复制带 Mutex 的结构体。`go vet` 能发现一部分。
- 在持有 `RWMutex` 读锁时再 `Lock`。

## 练习

用 `atomic.Int64` 重写一个只做计数的热路径，再用 Mutex 保护一个「计数 + 最近一次错误字符串」的结构，说明为什么后者不能拆成两个原子变量。
