# 专题 02：向量数据库（Milvus / PGVector）与 RAG 混合召回重排链路

## 📌 理论精要与现代化 RAG 架构

大语言模型存在“幻觉（Hallucination）”与“训练数据时效性滞后”两大天然硬伤。**检索增强生成（RAG - Retrieval-Augmented Generation）** 已成为企业级 AI 应用落地的核心基础设施：
1. **ANN 近似最近邻检索（Approximate Nearest Neighbor）**：针对高维向量空间（1536 维甚至更长），传统暴力暴算（Flat）在百万级数据下耗时不可接受。工业界采用 **HNSW（分层导航小世界图）** 或 **IVF-FLAT（倒排文件索引）** 实现毫秒级语义召回。
2. **混合召回（Hybrid Search）的必要性**：
   - **稠密向量（Dense Retrieval）**：擅长理解同义词、意图与深层语义泛化；但在专有名词、产品型号、工单编码、缩写词检索上极其容易漏召回；
   - **稀疏检索（Sparse Retrieval - BM25）**：擅长字面精确匹配，弥补向量模型对特定专业术语不敏感的短板。
3. **倒数排名融合（RRF - Reciprocal Rank Fusion）**：将来自向量空间和字面检索的异构归一化分数统一平滑映射（$Score = \sum \frac{1}{k + rank}$），避免分值尺度差异引发的偏差。

---

## 🏗️ 工业级 RAG 检索生成链路

```text
[用户提问] "解释 Go 调度器"
      │
      ├───────────────────────────────┐
      ▼ (Embedding 向量化)            ▼ (分词提取)
  [0.88, 0.12, 0.22...]          ["Go", "调度器"]
      │                               │
      ▼                               ▼
 [向量库 Milvus/PGVector]      [全文检索 Elastic/BM25]
 (HNSW 稠密召回 Top-K)         (字面稀疏召回 Top-K)
      │                               │
      └───────────────┬───────────────┘
                      ▼
            [ RRF 混合重排融合器 ]
                      │ (合并排序取出 Top-3 核心上下文)
                      ▼
         [ 提示词动态工程上下文构建 ]
                      │
                      ▼
         [ 送入 LLM 生成事实对齐的回答 ]
```

---

## 🔬 本机实测运行输出与基准数据

基于苹果 M1 Pro 芯片实测：

```bash
$ go test -v -bench=. ./04-ai-native-backend/02-vector-rag-pipeline/...
=== RUN   TestCosineSimilarity
--- PASS: TestCosineSimilarity (0.00s)
=== RUN   TestRAGPipelineRetrieval
--- PASS: TestRAGPipelineRetrieval (0.00s)
goos: darwin
goarch: arm64
pkg: lp-go/04-ai-native-backend/02-vector-rag-pipeline
cpu: Apple M1 Pro
BenchmarkCosineSimilarity-8       315496        3774 ns/op (1536维全浮点余弦相似度)
BenchmarkHybridSearchWithRRF-8    110637       10720 ns/op (百量级混合召回与RRF重排仅10微秒)
PASS
```

### 💡 核心结论与工程亮点：
- **1536 维超高速纯 Go 向量计算**：在 M1 Pro 上单次 1536 维浮点余弦相似度计算仅耗时 **3.77 微秒**。
- **高并发混合重排**：端到端百级文档向量 + 词频检索 + RRF 融合仅需 **10.72 微秒**，比网络 I/O 耗时低 3 个数量级。
