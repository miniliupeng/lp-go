# 专题 03：微服务可观测性（OpenTelemetry 链路贯穿 + Prometheus RED 模型）

## 📌 理论精要与现代化架构

微服务架构错综复杂的网络调用下，单点排查几乎不可行。现代可观测性体系（Observability）由三大支柱演进为统一的 OTel 标准：
1. **OpenTelemetry (OTel) 链路追踪**：遵循 W3C Trace Context 规范（`traceparent: 00-{trace_id}-{span_id}-{flags}`），通过在 `context.Context` 与跨网络协议头（HTTP Header / gRPC Metadata）中注入和提取，实现一条业务请求跨越数十个微服务的全局血缘追踪。
2. **Prometheus RED 指标模型**：微服务请求监控黄金三要素：
   - **Rate（吞吐量）**：每秒处理请求数。
   - **Errors（错误率）**：失败或异常请求占比。
   - **Duration（时延分布）**：耗时直方图分布（重点盯防 P99/P999 尾延迟）。

---

## 🏗️ 链路贯穿数据流向

```text
[HTTP Gateway]
    │ (生成 TraceID=abcd, SpanID=0001)
    ▼
    ├─ Injects "traceparent: 00-abcd-0001-01"
    ▼
[Order Service]
    │ (解析 TraceID=abcd, 派生子 SpanID=0002, 记录 RED 监控)
    ▼
    ├─ Injects "traceparent: 00-abcd-0002-01"
    ▼
[Payment Service]
      (维持 TraceID=abcd, 派生子 SpanID=0003, 完成支付上报)
```

---

## 🔬 本机实测运行输出与基准数据

基于苹果 M1 Pro 实测：

```bash
$ go test -v -bench=. ./03-microservice-governance/03-observability/...
=== RUN   TestTraceContextPropagation
--- PASS: TestTraceContextPropagation (0.00s)
=== RUN   TestREDMetrics
--- PASS: TestREDMetrics (0.00s)
goos: darwin
goarch: arm64
pkg: lp-go/03-microservice-governance/03-observability
cpu: Apple M1 Pro
BenchmarkTraceInjectExtract-8        762645          1562 ns/op
BenchmarkREDMetricsRecord-8        57246901            20.94 ns/op
PASS
```

### 💡 性能亮点：
- **极低开销指标记录**：基于内存 CAS 与分桶累加，单次 Prometheus RED 指标打点仅耗时 **20.94ns**，在数万 QPS 下对业务时延近乎零侵入。
- **W3C 标准协议规范化**：严格遵循分布式追踪工业标准，与 Jaeger、Zipkin、SkyWalking 和 Datadog 具备 100% 互操作性。
