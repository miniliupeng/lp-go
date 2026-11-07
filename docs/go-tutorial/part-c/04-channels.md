# 04 channel 与 select

源码依据：Go 1.27.1 `runtime/chan.go` 里 `hchan` 的字段和不变量。语言行为以规范为准。

## 核心原理

`hchan` 里和面试相关的字段有：`qcount`（缓冲区里现有的元素数）、`dataqsiz`（缓冲容量）、`buf`、`closed`、`sendx` / `recvx`、`recvq` / `sendq`，以及保护这些字段的 `lock`。等待者是 `sudog`。

`chan.go` 写明的不变量：

- `sendq` 和 `recvq` 至少有一个为空。例外是：无缓冲 channel 上，同一个 goroutine 在一个 `select` 里同时等发送和接收。
- 对有缓冲 channel：`qcount > 0` 说明 `recvq` 为空；`qcount < dataqsiz` 说明 `sendq` 为空。有空位就不会有人堵在发送，有数据就不会有人堵在接收。

无缓冲 channel 是一次会合。有缓冲 channel 在未满时发送方只把值拷进 `buf`。

## 关闭规则

这些是语言规则，示例在 [chan_test.go](../examples/partc/chanx/chan_test.go) 里都有测试：

| 操作 | 结果 |
| --- | --- |
| 从已关闭且已取空的 channel 接收 | 零值，`ok == false` |
| 关闭前缓冲里还有值 | 先把缓冲取完，`ok` 仍为 true |
| 向已关闭的 channel 发送 | panic |
| `close` 已关闭的 channel，或 `close(nil)` | panic |
| 对 nil channel 发送或接收 | 永远阻塞 |
| `for range` 已关闭的 channel | 取完剩余元素后结束 |

只由发送方关闭，或者由唯一负责关闭的那一方关闭。多个发送方时，用一个额外的信号 channel 或 `context` 通知它们退出，不要让每个发送方都 `close` 数据 channel。

## select

规范：如果有多个通信可以立即进行，实现用均匀的伪随机选一个。所以两个 `case` 都就绪时，测试不能断言一定走哪一个，只能断言走到的那个值是对的。

没有任何 case 就绪时，有 `default` 就走 `default`。nil channel 的 case 永远不会被选中，可以用来动态关掉一条分支：把那个 channel 变量设成 nil。

## 计时器通道

`time.NewTimer` 的文档写明：Go 1.23 起，计时器通道是同步的（容量 0）。Go 1.27 删除了 `asynctimerchan` 这个 GODEBUG，不能再退回旧的缓冲行为。在本仓库的 Go 1.27.1 上，`cap(time.NewTimer(d).C)` 和 `Ticker` 的 `C` 都是 0。`Stop` / `Reset` 之后不会再读到旧时间。未被引用的 Timer 和 Ticker 可以被 GC 回收，不必为了「防泄漏」而 `Stop`；要提前取消时仍然调用 `Stop`。

## 泄漏

G 堵在一个再也等不到对方的 channel 上，就是泄漏。常见原因：接收方已经退出却没有人再接收；向 nil channel 发送；`WaitGroup` 和 channel 互相等待。Go 1.27 的 `goroutineleak` profile 能发现「堵在同步原语上，而且这个原语从任何可运行的 G 都到不了」的那一类。它发现不了通过全局变量或仍可运行的 G 的局部变量还能摸到的阻塞。见 [性能一章](07-pprof-trace-bench.md)。

## 面试常问

- 画出无缓冲发送的双方谁先醒。
- 解释 `qcount > 0` 时为什么 `recvq` 必须为空。
- `select` 里两个 case 都就绪时能不能依赖顺序：不能。
- 为什么不要从接收方关闭：发送方还会发送，会 panic。

## 对比陷阱

- 「关闭 channel 会让发送方收到通知」。关闭不是广播给发送方的错误值；发送会 panic。
- 「缓冲 channel 的发送一定不阻塞」。缓冲满了就会阻塞。
- 把 `time.After` 放进高频循环。每次都新建计时器。Go 1.23 以后未引用的计时器可以被回收，但循环里仍然在分配。复用 `Timer` 并 `Reset`。

## 运行

```bash
go test ./docs/go-tutorial/examples/partc/chanx -count=1
```

## 练习

写一个 `select`：数据 channel 和 `ctx.Done()`。`ctx` 已取消时，循环必须退出，不能再依赖数据 channel。
