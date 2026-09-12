# 阶段六: AI 原生后端系统实战

2026 年的后端工程师必须掌握与大模型及向量数据库原生集成的能力。本阶段聚焦 AI 原生后端的三大支柱：HTTP SSE 长连接流式输出网关、高性能高维向量检索 RAG 链路，以及字节跳动开源的大模型应用编排框架 Eino。

---

## 💡 核心工程心智与底层机制

1. **SSE（Server-Sent Events）流式网关与反压机制**：
   - 基于 HTTP/1.1 长连接持续推送 `text/event-stream` 事件流；
   - 掌握通过 `http.Flusher` 实现毫秒级逐字打字机推流；
   - 核心防崩设计：构建基于缓冲 Channel 的反压流控（Backpressure Control），当客户端接收速率过慢或网络丢包时避免服务端 Goroutine 内存 OOM，支持客户端断开优雅级联取消（`r.Context().Done()`）。
2. **高维向量数据库与 RAG 检索链路**：
   - 深入稠密高维向量（Dense Embedding Vector）空间几何特征；
   - 掌握向量点积、模长计算与余弦相似度（Cosine Similarity）算法原理；
   - 掌握结合 Milvus / PGVector 搭建知识切片离线索引、在线相似度召回与 Prompt 动态拼装的完整 RAG 管道。
3. **Go 原生 AI 智能体应用编排（字节 Eino 框架）**：
   - 深入字节跳动专为 Go 语言打造的高性能大模型应用框架 **Eino**；
   - 掌握通过 JSON Schema 定义与大模型工具调用（Tool Calling / Function Calling）；
   - 掌握基于有向图（Graph）的状态机流转编排，构建具备自主规划、工具路由与上下文追踪能力的生产级 Agent。

---

## 🗺️ 阶段专题导航

| 专题目录 | 核心原语 | 关键掌握目标 |
| :--- | :--- | :--- |
| **[01-sse-streaming-gateway](./01-sse-streaming-gateway/)** | SSE 流式推送与反压防 OOM | 掌握 http.Flusher 实时打字机、Channel 背压缓冲与客户端断连熔断 |
| **[02-vector-rag-pipeline](./02-vector-rag-pipeline/)** | 向量空间检索与 RAG 管道 | 掌握余弦相似度比对、Top-K 高维召回与大模型上下文知识注入链路 |
| **[03-eino-agent](./03-eino-agent/)** | 字节 Eino 框架与 Agent 编排 | 掌握 Go 原生大模型工具调用、状态图编排与企业级智能体研发 |

---

## 🚀 统一运行验证

```bash
# 阶段全量单测与 AI 原语串行验证
go test -v ./06-ai-native-backend/...
```
