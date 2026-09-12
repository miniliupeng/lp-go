# 阶段五: 微服务架构与系统治理

在分布式与微服务架构中，单点可用性不再等于系统整体可用性。本阶段覆盖高性能 RPC 通信选型、生产级高可用防御治理（限流/熔断/降级/幂等），以及基于 OpenTelemetry 与 Prometheus 的全链路可观测性基建。

---

## 💡 核心工程心智与底层机制

1. **RPC 协议与底层传输引擎选型**：
   - 深入 gRPC 的 HTTP/2 二进制分帧、多路复用与 Protobuf 序列化协议；
   - 掌握字节跳动开源的 Kitex 高性能框架与自研 `Netpoll` 反应堆 I/O 模型（规避原生 `net.Conn` 每个连接一个协程的内存开销，实现跨连接缓冲区复用）。
2. **微服务高可用韧性防御治理**：
   - **限流（Rate Limiting）**：基于 `golang.org/x/time/rate` 令牌桶算法，平滑处理突发流量与削峰填谷；
   - **熔断（Circuit Breaker）**：基于 Sony `gobreaker` 状态机（Closed $\to$ Open $\to$ Half-Open），实现下游故障服务快速熔断与探针自愈；
   - **幂等（Idempotency）**：全局唯一流水号与业务前置 Token 验证，防范重试雪崩与接口重复提交。
3. **现代全链路可观测性体系（Observability）**：
   - 彻底贯彻云原生 **RED 指标模型**（Rate 吞吐、Errors 错误率、Duration 耗时分位数 P95/P99）；
   - 基于 **OpenTelemetry (OTel)** 规范实施跨服务链路注入与提取（W3C Trace Context 标准透传）；
   - 掌握 Prometheus Metric 打点、分布式 Span 传播与 Grafana 统一看板集成。

---

## 🗺️ 阶段专题导航

| 专题目录 | 核心原语 | 关键掌握目标 |
| :--- | :--- | :--- |
| **[01-rpc-kitex-grpc](./01-rpc-kitex-grpc/)** | gRPC 与 Kitex Netpoll 性能 | 深入 Protobuf 契约、gRPC 多路复用与 Kitex 自研网络模型优化 |
| **[02-reliability](./02-reliability/)** | 限流、断路器与接口幂等治理 | 掌握令牌桶动态流量塑形、熔断状态机降级与请求幂等防重实战 |
| **[03-observability](./03-observability/)** | OpenTelemetry 链路与 RED 监控 | 掌握 OTel W3C 上下文透传、Span 埋点与 Prometheus 延迟分位数监控 |

---

## 🚀 统一运行验证

```bash
# 阶段全量单测与治理模拟验证
go test -v ./05-microservice-governance/...
```
