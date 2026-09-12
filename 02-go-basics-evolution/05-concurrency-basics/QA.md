# 专题 05：Channel 状态矩阵与并发控制大厂高频面试题

---

### Q1: 请背诵并解析 Channel 的“3 态 5 操”矩阵（大厂硬核必考）？

**标准回答**：
- **三种状态**：
  1. `nil`（声明但未 `make` 的通道）
  2. `Open`（正常打开已初始化的通道）
  3. `Closed`（已调用 `close()` 关闭的通道）
- **关键操作表现**：
  - **向 nil 读/写**：**永久阻塞**。当前 Goroutine 挂起，若无其他可运行协程则引发致命死锁 `fatal error: all goroutines are asleep - deadlock!`。
  - **向 closed 写**：**立即触发 Panic**（`panic: send on closed channel`）。
  - **关闭 nil 或重复关闭 closed**：**立即触发 Panic**（`panic: close of nil channel` 或 `panic: close of closed channel`）。
  - **从 closed 读**：只要缓冲区还有数据，优先读取旧数据；一旦排空，**立即返回元素类型零值，且第二个布尔标志位 ok 为 false**，绝不阻塞。

---

### Q2: 为什么向无缓冲 Channel 发送数据，如果没有其他协程接收会发生死锁？

**标准回答**：
- **同步交接语义**：无缓冲 Channel 的容量为 0，发送方与接收方必须在时间线上**同时完成同步握手**。
- **死锁场景**：若在单协程顺序执行代码中直接写 `ch <- 1`，由于当前主线程正在等待接收者出现，而接收者代码还在当前行的后方，程序根本无法向下执行到接收逻辑，导致该协程自我阻塞，被运行时死锁检测器直接判定为全死锁（Deadlock）。
- **解决标准**：发送或接收的其中一方**必须由独立的 Goroutine 异步执行**。

---

### Q3: `select` 的多分支执行顺序是怎样的？如果多个 case 同时就绪会发生什么？

**标准回答**：
- **机制原理**：Go 的 `select` 借鉴了 Dijkstra 的守护命令（Guarded Commands）思想。
- **伪随机轮询**：当有多个 `case` 表达式同时满足就绪条件时，Go 运行时**通过伪随机数生成器（`fastrandn`）随机选择其中一个分支执行**！
- **设计考量**：如果像传统 if-else 一样按顺序自上而下匹配，排在最上面的高频通道会永久霸占执行机会，导致下方通道严重饿死（Starvation）。随机化保证了每个通信通道的公平性。

---

### Q4: `context.Context` 的四种核心派生函数有什么区别与使用场景？

**标准回答**：
1. **`context.WithCancel(parent)`**：返回子 context 与 `cancel` 函数。适用于由发起方主动触发取消信号（如任务提前完成、错误中断）。
2. **`context.WithTimeout(parent, duration)`**：设置相对超时间隔（如 500ms），时间到达自动触发取消。常用于下游 RPC 调用、数据库查询保护。
3. **`context.WithDeadline(parent, time)`**：设置绝对时间截止点（如具体的某个时间戳）。
4. **`context.WithValue(parent, key, val)`**：传递不可变的请求级别元数据（如 `TraceID`、用户鉴权身份信息）。注意：**严禁用 Value 传递可选的业务参数**。
