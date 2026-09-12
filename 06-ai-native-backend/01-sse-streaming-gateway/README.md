# 专题 01：大模型流式输出 SSE 网关、反压机制与长连接保活

## 📌 理论精要与 AI 时代网关架构

在 GenAI 大模型时代（2026 后端标配），传统请求-响应模式已无法满足大模型几十秒的生成延迟。**Server-Sent Events (SSE)** 凭借单向基于 HTTP 流式传输、天然穿透防火墙、自动重连与轻量特性，成为流式交互的绝对首选。

然而，AI 网关落地面临三大工程挑战：
1. **网络反压（Backpressure）**：当客户端网络较慢（如弱网手机）未能及时接收 Token，网关若无界缓冲上游推理结果，极易导致内存撑爆 OOM。必须引入有界 Channel 与超时丢弃/暂停反压。
2. **长连接保活（Keep-Alive）与中间代理超时**：Nginx、ALB 等代理若 60s 内无数据流动会主动切断连接。SSE 服务端必须按固定周期（如 15s）发送 `: ping\n\n` 注释心跳帧。
3. **即时刷新（Flusher）**：Go 标准库 `http.ResponseWriter` 默认带缓冲，必须主动断言为 `http.Flusher` 并在每个 Token 写入后执行 `Flush()`，确保首字时延（TTFT - Time to First Token）最小化。

---

## 🏗️ 流式网关与反压模型

```text
[LLM 推理上游]
      │ (持续吐出 Token: 10~50ms/token)
      ▼
┌────────────────────────────────────────────────────────┐
│              Go 高并发 SSE 网关 (有界管道)               │
│                                                        │
│   Push(Token) ──> [ 有界缓冲 Channel (容量 N) ]        │
│                         │ (下游消费慢? 触发反压超时拒绝)   │
│                         ▼                              │
│   Flusher.Flush() ──> HTTP/1.1 or HTTP/2 传输帧        │
│   Ticker ──> 定时 `: ping\n\n` 心跳注入                 │
└────────────────────────────────────────────────────────┘
      │
      ▼ (标准 text/event-stream)
[客户端 Web / App]
```

---

## 🔬 本机实测运行输出与基准数据

基于苹果 M1 Pro 实测：

```bash
$ go test -v -bench=. ./04-ai-native-backend/01-sse-streaming-gateway/...
=== RUN   TestSSEGatewayHandler
--- PASS: TestSSEGatewayHandler (0.13s)
=== RUN   TestBoundedStreamChannelBackpressure
--- PASS: TestBoundedStreamChannelBackpressure (0.02s)
goos: darwin
goarch: arm64
pkg: lp-go/04-ai-native-backend/01-sse-streaming-gateway
cpu: Apple M1 Pro
BenchmarkFormatSSE-8         3497174           339.7 ns/op
PASS
```

### 💡 核心结论与工程亮点：
- **反压保护实证**：在单测 `TestBoundedStreamChannelBackpressure` 中，当下游管道阻塞超阈值时，网关立即触发反压机制并切断无效堆积，成功保障网关内存安全。
- **SSE 帧极速序列化**：单次 SSE 标准文本协议拼装与 JSON 封装仅耗时 **339.7ns**，单机单核可轻松支持万级并发流式转发。
