# 专题 01：微服务 RPC 核心架构与编解码高频面试题 (QA)

### Q1：为什么字节跳动要自研 Kitex 而不直接使用官方 gRPC？
**标准深度解析**：
1. **网络模型瓶颈**：官方 gRPC 基于 Go 标准库 `net.Conn`，采用阻塞 I/O 模型，每个 TCP 连接需要绑定至少 2 个协程（一个读协程一个写协程）。在微服务网状拓扑、成千上万长连接下，数万协程造成巨大的上下文切换开销和内存占用。Kitex 引入底层 **Netpoll** 库，采用 epoll/kqueue 反应堆模型，用极少量的事件循环协程管理海量连接。
2. **内存拷贝与 GC 开销**：标准库 I/O 会将数据从内核态复制到用户态 buffer，造成频繁的堆分配。Kitex 基于 `LinkBuffer` 提供流式内存池，配合零拷贝 API，极大减少堆分配与 GC 压力。
3. **编解码序列化瓶颈**：gRPC Protobuf 依赖反射或频繁小指针对象分配，Kitex 的 FastCodec / Sonic 等编解码器通过预分配平坦内存结构和汇编加速，性能提高 2~3 倍。

---

### Q2：gRPC 是如何利用 HTTP/2 实现多路复用的？为什么还会有队头阻塞问题？
**标准深度解析**：
1. **HTTP/2 多路复用原理**：单条 TCP 连接上通过划分 Stream ID，将不同的 RPC 请求帧（HEADERS Frame / DATA Frame）交织并发传输，客户端与服务端按 Stream ID 分发重组，无需像 HTTP/1.1 那样排队阻塞。
2. **传输层队头阻塞（TCP HOL Blocking）**：HTTP/2 仅解决了应用层 HTTP 报文的队头阻塞。但由于底层依然是单条 TCP 连接，若网络发生丢包，TCP 协议为保证有序交付，内核滑动窗口必须等待丢包重传，导致后续所有 Stream 的数据包均被阻塞在内核缓冲区中。这也是为什么 HTTP/3 转向基于 UDP 的 QUIC 协议。

---

### Q3：微服务 RPC 调用中，如何优雅设计客户端的超时传递（Timeout Propagation）？
**标准深度解析**：
1. **Context Deadline 透传**：在发起 RPC 调用时，利用 `context.WithTimeout` 设置绝对截止时间 `Deadline`。
2. **跨进程协议头注入**：客户端中间件获取 `ctx.Deadline()`，将其换算为剩余毫秒数，写入 RPC 协议元数据（如 gRPC metadata `grpc-timeout`）。
3. **下游服务端拦截与约束**：下游服务拦截器收到请求后，解析该 Header，使用 `context.WithDeadline` 或 `min(本地默认超时, 上游透传剩余时间)` 生成新的 Context。如果链路上游因超时早已放弃，下游能够立即感知并快速失败，防止雪崩与无效计算。
