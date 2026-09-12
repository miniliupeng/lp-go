# 专题 02：微服务高可用三板斧（限流、熔断降级、业务幂等）

## 📌 理论精要与生产落地

在分布式高并发微服务集群中，高可用治理是保障系统不发生雪崩的三道防火墙：
1. **令牌桶限流（Token Bucket）**：应对突发流量（Bursty Traffic）的最佳实践。通过匀速向桶中补充令牌，允许短时间突发请求消耗桶内存量令牌，兼顾吞吐量与平均速率控制。
2. **自适应熔断（Circuit Breaker）**：微服务网络调用天然存在抖动。熔断器（Closed $\leftrightarrow$ Open $\leftrightarrow$ Half-Open 状态机）在下游依赖故障时及时“拉闸断电”，避免线程/协程级联挂起导致全链路雪崩，并在超时后通过半开探测实现自动故障愈合。
3. **业务幂等性 Token 机制**：面对弱网重试、用户连击、消息重投，通过唯一幂等 Key（如雪花算法生成）配合原子并发锁与结果暂存，实现“至少一次传输”下的“恰好一次处理”。

---

## 🏗️ 熔断器状态机模型

```text
       ┌───────────────┐
       │    CLOSED     │ ──(错误数超标)──> ┌───────────────┐
       │ (正常放行请求) │                  │     OPEN      │
       └───────────────┘                  │ (快速失败拦截)│
               ▲                          └───────────────┘
               │                                  │
         (连续探测成功)                      (超时冷却到期)
               │                                  ▼
       ┌────────────────┐                  ┌───────────────┐
       │   HALF-OPEN    │ <────────────────│   HALF-OPEN   │
       │ (放行探测流量) │                  │ (进入半开探测)│
       └────────────────┘                  └───────────────┘
```

---

## 🔬 本机实测运行输出与基准数据

基于苹果 M1 Pro 实测指标：

```bash
$ go test -v -bench=. ./03-microservice-governance/02-reliability/...
=== RUN   TestTokenBucketLimiter
--- PASS: TestTokenBucketLimiter (0.03s)
=== RUN   TestCircuitBreakerStateTransitions
--- PASS: TestCircuitBreakerStateTransitions (0.06s)
=== RUN   TestIdempotency
--- PASS: TestIdempotency (0.00s)
goos: darwin
goarch: arm64
pkg: lp-go/03-microservice-governance/02-reliability
cpu: Apple M1 Pro
BenchmarkTokenBucketLimiter-8       18667584         64.12 ns/op
BenchmarkIdempotencyManager-8       16093026         71.56 ns/op
PASS
```

### 💡 性能亮点：
- **无协程消耗的惰性令牌桶**：每次请求通过时钟差值 `elapsed * rate` 纯数学补算，单次判定开销仅 **64ns**，无需后台后台开启常驻协程打点计时。
- **并发幂等原子判定**：基于 `sync.Map` 的 CAS/LoadOrStore 思想，单次幂等比对判定耗时 **71ns**，极具生产实用价值。
