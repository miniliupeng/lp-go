# 专题 05：pprof 性能剖析与高并发线上问题现场实战大厂高频面试题

---

### Q1: 生产环境中，微服务如何安全、低开销地集成 pprof 探针？

**标准回答**：
1. **标准集成**：直接通过空白导入 `import _ "net/http/pprof"`，并启动一个**独立的内部监控端口**（如仅监听私有网络端口 `localhost:6060` 或容器内端口），**严禁将 pprof 路由直接暴露在外网公网网关上**，防止攻击者通过反复采集导致 DoS 或泄露内部内存数据结构。
2. **性能开销**：
   - 默认情况下（未发起采样 HTTP 请求时），pprof 的 CPU 采集器处于**休眠关闭状态**，无额外开销；
   - 内存采样（Heap Profile）依托内部的 `mallocgc`，开销微乎其微（< 1%）；
   - 仅当管理员或自动化系统通过 HTTP 访问 `/debug/pprof/profile?seconds=30` 时，内核才会激活 `SIGPROF` 信号进行 30 秒采样，开销约为 1%~3% CPU，对高可用服务完全安全。

---

### Q2: 线上发生 OOM 时，为什么看 `inuse_space` 有时看不出问题，应该怎么比对分析？

**标准回答**：
- **`inuse_space` vs `alloc_space`**：
  - `inuse_space`：记录的是**当前时刻依然存活未被 GC 的堆内存**；
  - `alloc_space`：记录的是自程序启动以来**累计分配的总内存量**（包括已经被 GC 回收的内存）；
- **排查策略**：
  1. **慢速持续泄漏（Slow Leak）**：看 `inuse_space`。通常是由于全局 Map 累积、切片截取大数组无法释放，导致当前存活内存稳步攀升。
  2. **高频临时垃圾引发的内存抖动（High-frequency Alloc Spikes）**：如果某个接口瞬间大量并发，单次调用生成巨大的短生命周期对象（如 JSON 解析几百 MB 的文本），即使 GC 能回收，但在瞬时突发流量下容易瞬间突破容器 cgroup 限额触发 OOM。此时 `inuse_space` 看起来很小，但切换到 `alloc_space` 或 `alloc_objects` 会立即发现该函数存在恐怖的累计分配量。
  3. **差异比对（Diff Analysis）**：在相隔 10 分钟后分别采集两份堆画像，使用 `go tool pprof -base mem_base.pprof mem_latest.pprof`，直观查看净增量。

---

### Q3: 生产环境中如何精准排查 Goroutine 泄漏？

**标准回答**：
- **第一步：监控告警**：通过 Prometheus 暴露的 `go_goroutines` 指标，发现 Goroutine 数量呈阶梯式单调上升，且流量低峰期依然不下降。
- **第二步：采集分析**：
  - 通过 curl 导出协程画像：`curl http://localhost:6060/debug/pprof/goroutine?debug=2 -o goroutine.txt`（`debug=2` 会将每个协程的完整函数调用堆栈以文本形式可读打印）；
  - 或通过 `go tool pprof http://localhost:6060/debug/pprof/goroutine` 执行 `top`。
- **第三步：识别特征模式**：
  - 观察排在最前列、数量高达数千上万个的协程堆栈，如果它们全部停留在 `chan send`、`chan receive`、`sync.WaitGroup.Wait` 或 `net.Conn.Read`，即可瞬间定位阻塞的代码行。

---

### Q4: 什么是火焰图中的“平顶山”？如何根据平顶山制定优化策略？

**标准回答**：
- **平顶山特征**：在 CPU 火焰图中，纵轴是调用栈，横轴是抽样占比。如果某个函数条块位于火焰图的最顶端（没有子函数），并且其水平跨度非常宽，呈现出像平顶山一样的形状，说明**该函数自身内部包含了极其耗时的计算逻辑**。
- **优化决策**：
  1. 如果平顶山是**业务逻辑函数**（如高频循环、正则、加密计算）：对算法进行优化，避免不必要的重复计算，引入缓存或改用无内存逃逸的算法；
  2. 如果平顶山是 `runtime.mallocgc`：说明程序中存在极其严重的堆分配开销，需要引入 `sync.Pool` 对象复用，或通过切片预分配消除 `append` 扩容；
  3. 如果平顶山是 `runtime.futex` 或 `runtime.sync`：说明程序存在严重的并发锁争用，需要拆小临界区、采用读写锁、无锁队列（Atomic）或分段锁。
