# 专题 03：微服务可观测性高频面试题 (QA)

### Q1：为什么分布式链路追踪要统一收敛到 W3C Trace Context 规范？
**标准深度解析**：
1. **打破供应商绑定与协议分裂**：过去业界存在 Zipkin（`X-B3-TraceId`）、Jaeger（`uber-trace-id`）、SkyWalking 等多种异构 Header。跨团队、跨多语言网关或跨云服务调用时，链路极易中断或需要维护极其臃肿的 Header 转换中间件。
2. **W3C 规范的标准性**：统一采用 `traceparent: 00-{trace_id}-{span_id}-{flags}`，包含版本号、全局 16 字节追踪 ID、8 字节父跨度 ID 及采样标记。此外还提供 `tracestate` 用于透传供应商私有状态。现代 OpenTelemetry SDK 原生全面支持 W3C 格式。

---

### Q2：Prometheus RED 模型与 Google SRE 四个黄金指标有何异同？
**标准深度解析**：
1. **Google SRE 四个黄金指标**：延迟（Latency）、流量（Traffic）、错误（Errors）、饱和度（Saturation）。
2. **RED 模型（Tom Wilkie 提出）**：专为请求驱动的微服务架构简化提炼而来：
   - **Rate**（对应 Traffic）：每秒请求速率；
   - **Errors**（对应 Errors）：失败请求数；
   - **Duration**（对应 Latency）：请求耗时分布。
3. **区别与关系**：RED 模型聚焦于服务本身对客户端承诺的 SLA 表现；而 Saturation（饱和度，如 CPU/内存/连接池使用率）属于资源维度指标（通常与 USE 模型结合），二者互为补充。

---

### Q3：高并发系统中，开启全量分布式链路追踪会导致什么灾难？生产环境如何治理？
**标准深度解析**：
1. **灾难隐患**：若每秒 10 万 QPS 全部生成完整 Span，跨网络上报给 OTel Collector / Elasticsearch 会占用数十倍于业务流量的网络带宽与存储 I/O，并可能因 GC 停顿剧烈拖垮业务核心服务。
2. **生产级治理策略（采样率控制）**：
   - **头部采样（Head-based Sampling）**：在网关或链路起点按照固定比例（如 1%）或哈希算法打上 Sampled 标记，后续链路只透传已采样的追踪；
   - **尾部采样（Tail-based Sampling）**：Collector 端缓存整个 Trace 的所有 Span，等待请求结束时判定：**所有报错请求 100% 采集，所有耗时大于 1 秒的慢请求 100% 采集**，普通正常请求仅低频抽样。这是性价比最高的生产实践。
