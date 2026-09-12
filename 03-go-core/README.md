# 阶段三: Go 核心与底层系统

Go 语言之所以在云原生与高性能微服务领域占据统治地位，其灵魂全在于 Go Runtime 底层的轻量调度、自动化内存治理与原生高并发并发原语。本阶段直击大厂 T0 淘汰线核心，破译 Runtime 黑盒。

---

## 💡 核心工程心智与底层机制

1. **GMP 调度器模型与协作式抢占**：
   - Goroutine（用户态轻量线程，2KB 起步）、Machine（操作系统内核线程 M）与 Processor（逻辑处理器 P，绑定本地运行队列 LRQ）；
   - 掌握 Work-stealing（工作窃取算法）与 Hand-off（P/M 分离转移）调度策略；
   - 深刻理解基于系统调用与 `sysmon` 信号驱动的非协作式抢占机制（Preemption）。
2. **三色标记清除算法与写屏障机制**：
   - 深入白色（未访问）、灰色（已访问但引用的子对象未扫描）、黑色（存活且子对象已扫描）三色标记状态机；
   - 掌握强三色不变性与弱三色不变性，破译插入写屏障与删除写屏障融合的**混合写屏障（Hybrid Write Barrier）**；
   - 生产环境合理配置 `GOMEMLIMIT` 与 `GOGC`，彻底治理软内存上限防 OOM。
3. **hchan 底层拓扑与同步原语**：
   - Channel 底层是由环形缓冲区（`buf`）、互斥锁（`lock`）、发送等待队列（`sendq`）与接收等待队列（`recvq`）构成的 `hchan` 结构体；
   - 深刻理解 `sync.Mutex` 从正常模式到饥饿模式（8 个调度周期未获取锁强制饥饿抢占）的自适应自旋切换；
   - 掌握 `sync.Pool` 双向链表与私有对象池，大幅减免 GC STW 标记开销。
4. **Go 现代运行时与编译优化**：
   - Go 1.23+ 原生 `iter.Seq` 迭代器模式解耦流式处理；
   - 现代结构化高性能日志库 `log/slog` 生产选型；
   - 利用 Profile-Guided Optimization (PGO) 基于生产真实采样驱动编译器精准内联提升 5%~15% 吞吐。
5. **Pprof 生产线上性能调优 SOP**：
   - 掌握 CPU Profiling、Heap 内存泄漏定位、Goroutine 悬挂阻塞与 Mutex 锁争用现场复盘与诊断。

---

## 🗺️ 阶段专题导航

| 专题目录 | 核心原语 | 关键掌握目标 |
| :--- | :--- | :--- |
| **[01-gmp-scheduler](./01-gmp-scheduler/)** | GMP 调度器与协作抢占 | 掌握 GMP 模型拓扑、Work-stealing、sysmon 抢占与 Trace 工具链实战 |
| **[02-memory-gc](./02-memory-gc/)** | 内存逃逸与三色并发 GC | 掌握堆栈逃逸分析、混合写屏障、STW 耗时与 GOMEMLIMIT 防 OOM 治理 |
| **[03-channel-sync](./03-channel-sync/)** | hchan 源码与 Mutex 饥饿机制 | 深入 Channel 底层环形缓冲区、Mutex 饥饿模式与 sync.Pool 对象池优化 |
| **[04-modern-runtime-features](./04-modern-runtime-features/)** | 现代 Runtime、iter 与 PGO | 掌握 Go 1.22+ 现代特性、iter 迭代器范式、slog 日志与 PGO 编译优化 |
| **[05-pprof-tuning](./05-pprof-tuning/)** | Pprof 生产性能压测排障 | 掌握 CPU/Heap/Goroutine/Block 剖析，实战定位线上死锁与内存泄漏 |

---

## 🚀 统一运行验证

```bash
# 阶段全量单测与基准测试串行验证
go test -v ./03-go-core/...
```
