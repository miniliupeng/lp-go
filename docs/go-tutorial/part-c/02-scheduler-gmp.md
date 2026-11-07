# 02 GMP 调度

源码依据：Go 1.27.1 的 `runtime/proc.go` 开头注释，以及 `runtime/preempt.go`、`runtime/signal_unix.go`。

## 核心原理

调度器把可运行的 goroutine 分到工作线程上。三个对象：

| 名字 | 是什么 |
| --- | --- |
| G | goroutine |
| M | 操作系统线程。执行 Go 代码时必须持有一个 P；阻塞在系统调用里时可以不持有 P |
| P | 执行 Go 代码所需的资源，数量大致等于 `GOMAXPROCS` |

设计文档在 <https://golang.org/s/go11sched>。状态是分散的，尤其是每个 P 自己的运行队列，因此快路径上算不出一个精确的全局谓词。

有可运行的 G 时，调度器不会为每个 G 都去唤醒一条线程。`proc.go` 写明的策略是：提交工作之后，只有在「存在空闲 P」并且「当前没有 spinning 线程」时才额外唤醒一条线程。spinning 指这条线程手里没有本地工作，也还没从全局队列或 netpoller 找到工作。它会再去别的 P 的队列、定时器堆和 GC 工作里找；找到就退出 spinning 并执行，找不到就退出 spinning 然后休眠。最后一条 spinning 线程如果找到了工作，必须再唤醒一条新的 spinning 线程，否则可能长时间用不满 CPU。

网络等待走 netpoller：G 阻塞在网络上时，M 不必一直占着 P。这是「开很多连接不会开很多线程」的原因。阻塞发生在 poller 上，G 进入等待，P 可以去跑别的 G。数据到达后运行时再把 G 放回运行队列。

## G 的状态

`runtime/runtime2.go` 里和面试有关的状态（数字是 Go 1.27.1 源码里的 iota，不要背成永远不变的 ABI）：

| 状态 | 值 | 含义 |
| --- | --- | --- |
| `_Gidle` | 0 | 刚分配，还没初始化 |
| `_Grunnable` | 1 | 在某个运行队列上，没有在执行用户代码 |
| `_Grunning` | 2 | 可以执行用户代码，占用一个 M，通常还有 P。进出系统调用的很短窗口里可能暂时没有 P |
| `_Gsyscall` | 3 | 正在执行系统调用，占用 M，不拥有 P。这个状态里不能碰 `g.m.p` |
| `_Gwaiting` | 4 | 堵在运行时里，例如 channel。应当记在某个等待队列上，以便以后 `ready` |
| `_Gdead` | 6 | 已退出或在空闲链表上。5 是留给调试器的未使用值 |

另外两个不要和上面混在一起讲：`_Gpreempted`（9）是为了 `suspendG` 自己停下来，还没有人负责把它 `ready`；`_Gleaked`（10）表示 GC 认定这是泄漏的 goroutine。`_Gscan` 是 `0x1000`，和上面的状态按位组合，表示 GC 正在扫描这个栈。面试答到 runnable、running、syscall、waiting 这四个，并说 syscall 时不占着 P，就够了。

## 抢占

`runtime/preempt.go` 把安全点分成三类：

1. **阻塞安全点**：G 已经不在运行，堵在同步原语或系统调用里。
2. **同步安全点**：G 在函数序言里做栈界限检查时发现栈界被下毒，于是进入调度。
3. **异步安全点**：用户代码里可以安全暂停、并能保守扫描栈和寄存器的指令。运行时用操作系统机制把线程停住。Unix 上这个信号是 `SIGURG`（`signal_unix.go` 里 `sigPreempt = _SIGURG`），因为这个信号允许偶发、又很少被应用程序真正拿来用。

所以 Go 1.14 以后，一个完全不调用函数的死循环也可以被抢占。这和「只有函数调用才能让出 P」不是同一件事。

## 可运行示例说明

[sched 测试](../examples/partc/sched/preempt_test.go) 把 `GOMAXPROCS` 设为 1，再跑一个循环。循环里有 `atomic.Bool.Load`，那是函数调用，属于同步安全点。测试证明的是：只有一个 P 时，碰到安全点的 G 不会把另一个 G 饿死。它**没有**证明 `SIGURG`。异步抢占要看上面的源码，不要把这个测试的结果说成信号实验。

```bash
go test ./docs/go-tutorial/examples/partc/sched -count=1
```

测试结束会把 `GOMAXPROCS` 设回去。它是进程级设置，不要在库的 `init` 里改。

## 面试常问

- P 为什么存在：M 执行 Go 代码需要的本地缓存和运行队列。系统调用可以让出 P，别的 M 继续跑 Go 代码。
- 什么是 work stealing：自己的队列空了，从别的 P 偷 G。注释里的 spinning 是为了在「不要频繁唤醒线程」和「最终用满并行度」之间取平衡。
- 异步抢占用什么信号：Unix 上是 `SIGURG`。不要说「所有平台都是 SIGURG」。

## 对比陷阱

- G 不是线程。一个 M 可以先后跑很多 G。
- `GOMAXPROCS` 不是「最大线程数」。阻塞在系统调用里的 M 可以比 P 多。
- 协作式调度的老说法已经过时。同步安全点和异步安全点要分开讲。

## 练习

把测试里的 `atomic` 负载改成一次 `runtime.Gosched()`，说明这是主动让出，和信号抢占不是同一条路径。
