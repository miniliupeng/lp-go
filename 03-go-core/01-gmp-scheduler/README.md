# 专题 01：GMP 调度器模型与底层演进

## 1. 核心概念与关系架构

Go 的并发调度核心由三个实体协同完成：

| 实体 | 代表含义 | 职责与属性 |
| :--- | :--- | :--- |
| **G (Goroutine)** | 协程结构体 | 包含调用栈（初始 2KB，动态伸缩）、程序计数器（PC）、状态（`_Gidle`, `_Grunnable`, `_Grunning`, `_Gwaiting`, `_Gdead` 等）及关联的调度上下文。 |
| **M (Machine)** | 操作系统线程 | 由 OS 内核直接调度，数量上限默认为 10000（可通过 `debug.SetMaxThreads` 调整）。必须绑定 P 才能执行 Go 代码。 |
| **P (Processor)** | 逻辑处理器 | 虚拟的上下文资源池，数量默认为 CPU 核心数（`GOMAXPROCS`）。拥有本地运行队列（LRQ，容量 256），掌握执行 Go 代码所需的资源。 |

```
                全局运行队列 (GRQ)
               [ G1 -> G2 -> G3 ]
                       │
       ┌───────────────┴───────────────┐
       ▼                               ▼
    [ P1 ]                          [ P2 ]
本地队列 (LRQ: 256)             本地队列 (LRQ: 256)
 [ G4 -> G5 ]                    [ G6 -> G7 ]
       │                               │
       ▼                               ▼
    [ M1 ] (OS Thread)              [ M2 ] (OS Thread)
       │                               │
       ▼                               ▼
  CPU Core 1                      CPU Core 2
```

---

## 2. 调度循环与两级工作窃取（Work-Stealing）

当一个 M 绑定的 P 执行完当前 G 时，会触发 `schedule()` 调度循环寻找可运行的 G，寻找优先级如下：

1. **每 61 次循环**：优先检查一次**全局队列（GRQ）**，防止全局队列中的 G 饿死。
2. **本地队列（LRQ）**：从本 P 的运行队列头部弹出可执行 G。
3. **工作窃取（Work-stealing）**：若本地与全局皆为空，随机挑选另一个 P，从其 LRQ **尾部偷取一半（n/2）** 的 G。
4. **网络轮询器（Netpoller）**：检查是否有就绪的网络 I/O 事件相关的 G。

---

## 3. 阻塞处理机制：系统调用与 Hand-off

- **系统调用阻塞（Syscall Hand-off）**：
  - 当 M 执行阻塞的系统调用（如磁盘读写）时，M 进入内核态阻塞。
  - Go 监控线程 `sysmon` 发现此 P 长期处于 `_Psyscall` 状态（默认超过 10ms），会将 P 与该 M 解绑。
  - P 寻找空闲的 M（或新建 M）继续执行其他 G，保证 CPU 核心持续被利用。
- **系统调用退出**：
  - M 退出系统调用后，优先尝试拿回原来的 P；若失败则尝试获取全局空闲 P；若均失败，将 G 放回全局队列，M 进入休眠池。

---

## 4. 抢占式调度演进历史（面试高频考点）

- **Go 1.13 之前：协作式调度（Cooperative）**
  - 只能在特定点触发让出（如函数调用栈扩容检查、Channel 阻塞、主动调用 `runtime.Gosched()`）。
  - **严重缺陷**：如果 Goroutine 内是一个密集的纯数值计算死循环（无函数调用），该 G 会独占线程与 P，无法被中断，甚至导致 GC STW 停顿超长。
- **Go 1.14+ 至今：基于信号的非协作式抢占调度（Signal-based Preemption）**
  - `sysmon` 守护线程定期检测，发现某 G 运行超过 10ms 时，向绑定该 G 的 M 发送 `SIGURG` 信号。
  - OS 线程捕获 `SIGURG` 信号，进入信号处理函数 `sighandler()`。
  - 信号处理函数将目标 G 的 PC（程序计数器）重写为 `asyncPreempt` 函数，使 G 优雅保存现场并让出 CPU，彻底解决无函数调用死循环占用 P 的历史痛点。

---

## 5. 本机实测运行与基准数据（实操记录）

### ① 运行抢占实验与 Trace 导出
```bash
go run ./01-go-core/01-gmp-scheduler/main.go
```
**实测输出**：
```text
=== 1. GMP 基础环境信息 ===
逻辑 CPU 数量 (GOMAXPROCS): 8
当前活动 Goroutine 数量: 1

=== 2. 抢占式调度验证（单 P 场景）===
[G2 观察协程] 成功抢占执行！证明基于信号的异步抢占生效。
[G1 密集计算] 退出，累计计算轮次: 7550095
=== GMP 调度实验结束，已生成 trace.out ===
```
> **现象解读**：
> 在 `runtime.GOMAXPROCS(1)` 强制单核环境下，G1 是一个无任何系统调用和 I/O 的死循环。由于 Go 1.14+ 信号抢占机制，系统向执行线程注入了 `SIGURG` 信号，迫使 G1 让出 CPU，G2 仅延迟 10ms 就成功抢占并执行完毕，彻底避免了单核锁死。

### ② 调度与并发基准测试（Benchmark，采用现代 `b.Loop()` 规范）
```bash
go test -bench=. -benchmem -run=none ./01-go-core/01-gmp-scheduler/...
```
**实测数据（Apple M1 Pro / darwin-arm64）**：
```text
BenchmarkGoroutineCreation-8     3589736      327.3 ns/op      32 B/op      2 allocs/op
BenchmarkGoschedOverhead-8      16772728       72.64 ns/op      0 B/op      0 allocs/op
BenchmarkConcurrencyThroughput-8 1000000000     0.2147 ns/op     0 B/op      0 allocs/op
```
> **核心数据量化结论 & 现代 API 特性**：
> 1. **`b.Loop()` 优势（Go 1.24+ 引入）**：淘汰传统的 `for i := 0; i < b.N; i++` 样板代码，首轮迭代前自动完成 `b.ResetTimer()`，使初始化准备代码与核心基准代码完全隔离，同时防止编译器将无副作用的测量循环死代码消除。
> 2. **轻量协程开销**：单次创建并等待一个 Goroutine 仅耗时 **327.3 ns**，堆内存分配仅 **32 Bytes**（相比 OS 线程数 MB 内存与数微秒创建，差距达两个数量级）。
> 3. **主动出让成本极低**：调用 `runtime.Gosched()` 仅需 **72.6 ns**，且产生 **0 B** 堆内存分配。

### ③ 调度可视化分析（动手实操）
在当前目录下运行以下命令，即可启动本地追踪界面：
```bash
go tool trace trace.out
```
查看 `View trace`，可直观观测 Goroutine 在各核（Proc）上的状态迁移时间轴。

