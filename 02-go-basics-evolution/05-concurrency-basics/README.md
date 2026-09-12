# 专题 05：并发调度协同、Channel 状态矩阵与 Context 树级联穿透

> 本专题带领开发者系统跨入 Go 语言标志性的高并发世界：解构 `go func()` 协程生命周期与 `sync.WaitGroup` 协调规范、无缓冲与有缓冲通道通信模型、大厂必考淘汰线——**Channel 3态5操完整矩阵**、经典死锁边界、`select` 伪随机防饥饿机制，以及 `context.Context` 树形取消机制。

---

## 1. 协程生命周期协调：`go func()` 与 `sync.WaitGroup`

在任何函数或闭包前加上 `go` 关键字，就会立即由 Go 运行时启动一个轻量级协程（Goroutine）：

```go
// 启动子协程
go func(msg string) {
    fmt.Println(msg)
}("Hello from Goroutine!")
```

> [!CAUTION]
> **生产避坑**：严禁用 `time.Sleep` 盲目等待协程结束！生产环境必须使用 `sync.WaitGroup` 或 `Context` 进行确定性生命周期协调。
> - `wg.Add(1)` 必须在启动协程的外部调用；
> - `defer wg.Done()` 放置在协程内部第一行；
> - 主协程调用 `wg.Wait()` 安全阻塞。

---

## 2. 共享内存安全防线：互斥锁（`sync.Mutex`）与读写锁（`sync.RWMutex`）

虽然 Go 倡导通信共享内存，但在高并发本地缓存、全局计数器与状态机场景下，锁依然是性能最高、最直接的同步手段：

### 2.1 互斥锁（`sync.Mutex`）
提供完全排他的访问权限。任何协程在访问临界区（Critical Section）前调用 `mu.Lock()`，离开时通过 `defer mu.Unlock()` 释放：
```go
var mu sync.Mutex
var counter int

func Increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++ // 安全累加，绝无多协程数据撕裂（Data Race）
}
```

### 2.2 读写互斥锁（`sync.RWMutex`）
针对典型的“**读多写少**”（Read-heavy）服务场景：
- **读锁（`rw.RLock()` / `rw.RUnlock()`）**：允许多个协程同时持有读锁进行并发读取，互不阻塞；
- **写锁（`rw.Lock()` / `rw.Unlock()`）**：完全独占，写入时阻塞所有新的读锁与写锁。
- **性能优势**：在 90% 读、10% 写的缓存业务中，`RWMutex` 吞吐量通常比普通 `Mutex` 高出 **3~5 倍**。

---

## 3. 通道（Channel）通信模型：无缓冲 vs 有缓冲

Go 并发核心箴言：“**不要通过共享内存来通信，而要通过通信来共享内存**”。
- **无缓冲通道（`make(chan int)`）**：同步握手，发送方与接收方必须同时就绪；
- **有缓冲通道（`make(chan int, capacity)`）**：异步解耦，缓冲区未满发送不阻塞，缓冲区非空接收不阻塞。

---

## 4. 大厂核心考核淘汰线：Channel 3态 5操完整矩阵表

每一个大厂面试官必考的硬核基础，牢记下表：

| 操作 | 未初始化（`nil` Channel） | 正常打开（`Open`） | 已关闭（`Closed`） |
| :--- | :--- | :--- | :--- |
| **读（`<-ch`）** | **永久阻塞**（导致死锁） | 阻塞或成功读取数据 | **立即返回零值**，且 `ok == false` |
| **写（`ch <- x`）** | **永久阻塞**（导致死锁） | 阻塞或成功写入数据 | **直接触发 Panic**！💥 |
| **关（`close(ch)`）**| **直接触发 Panic**！💥 | 正常关闭 | **直接触发 Panic**（重复关闭）！💥 |
| **`len(ch)`** | 返回 `0` | 当前缓冲区元素个数 | 当前缓冲区剩余元素个数 |
| **`cap(ch)`** | 返回 `0` | 缓冲区总容量 | 缓冲区总容量 |

> **生产避坑铁律**：
> 1. **谁生产，谁关闭**：永远只由发送方（Producer）关闭 Channel，严禁由接收方关闭 Channel！
> 2. **多个并发生产者时不要关闭**：如果多个协程向同一个 Channel 写数据，不能由任何一个写协程关闭，可借助外部 `sync.Once` 或由外部协调者通过关闭 `stopCh` 信号通知退出。

---

## 5. `select` 多路复用的底层机理
- **伪随机轮询机制**：如果多个 `case` 通道同时就绪，Go 运行时**绝不按代码书写顺序执行**，而是通过伪随机算法（`fastrandn`）随机选择一个执行！这能彻底防止因书写顺序固定导致某个高频通道饥饿（Starvation）。
- **非阻塞读写（`default` 分支）**：当所有 `case` 均阻塞时，若存在 `default`，则立即走 `default`，实现非阻塞的尝试收发。

---

## 6. `context.Context` 树形取消与超时传播
在微服务链条中，一个 HTTP 请求可能衍生出数个下游 RPC 与数据库调用。若用户中途断开连接，必须级联通知所有子协程立即停止计算。

```
                context.Background() (根节点)
                          │
            WithTimeout(ctx, 100ms) (父节点)
                          │
          ┌───────────────┴───────────────┐
          ▼                               ▼
       Worker 1                        Worker 2
  (监听 <-ctx.Done())             (监听 <-ctx.Done())
```
- 一旦父 Context 超时或触发 `cancel()`，取消信号会顺着树形结构**瞬间向所有子孙 Context 广播**，各协程监听 `<-ctx.Done()` 即可实现毫秒级优雅退出。

---

## 7. 📊 本机实测运行与基准数据（实操记录）

### ① 运行并发与 Channel 状态验证
```bash
go run ./02-go-basics-evolution/05-concurrency-basics/main.go
```
**实测输出**：
```text
=== 1. 新手入门：sync.WaitGroup 优雅并发等待 ===
Worker 3 启动并执行完毕
Worker 1 启动并执行完毕
Worker 2 启动并执行完毕
所有 Worker 全部平稳完成！

=== 2. 基础入门：sync.Mutex 与 sync.RWMutex 并发安全保护 ===
互斥锁保护 1000 次并发累加结果: 1000 (预期: 1000)
读写锁 RLock 安全读取配置: turbofan

=== 3. 硬核底层：Channel 状态矩阵之从 Closed 读取 ===
关闭后第 1 次读取: val=100, ok=true (读取缓冲内数据)
关闭后第 2 次读取: val=200, ok=true (读取缓冲内数据)
关闭且排空后第 3 次读取: val=0 (返回类型零值), ok=false (明确告知已关闭)

=== 4. 硬核底层：select 多路复用与 default 非阻塞 ===
成功非阻塞写入 ping
通道已满，成功触发 default 非阻塞分支！

=== 5. 硬核底层：context.Context 树形超时级联取消 ===
Worker 2 接收到 Context 取消信号: context deadline exceeded，立即优雅退出！
Worker 1 接收到 Context 取消信号: context deadline exceeded，立即优雅退出！
```

### ② 无缓冲通道（同步）vs 有缓冲通道（异步）吞吐对比（Benchmark）
```bash
go test -v -bench=. -benchmem ./02-go-basics-evolution/05-concurrency-basics/...
```
**实测数据（Apple M1 Pro / darwin-arm64）**：
```text
=== RUN   TestClosedChannelRead
--- PASS: TestClosedChannelRead (0.00s)
BenchmarkUnbufferedChannel-8    8228481   130.00 ns/op   0 B/op   0 allocs/op
BenchmarkBufferedChannel-8     36519627    33.11 ns/op   0 B/op   0 allocs/op
PASS
```
> **量化结论**：
> 1. 单测严格验证：已关闭通道排空后读取立即返回零值与 `ok=false`，绝不阻塞。
> 2. **异步缓冲提升吞吐近 4 倍**：有缓冲通道单次收发仅需 **33.11 ns**，而无缓冲通道需要两个协程即时上下文握手（**130.0 ns**），在生产者/消费者速率不均场景下，有缓冲能极大地平滑抖动并提升吞吐。

