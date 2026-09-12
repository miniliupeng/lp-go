package main

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
)

// ============================================================================
// 1. 稠密向量相似度计算 (Cosine Similarity)
// ============================================================================

// CosineSimilarity 计算两个向量的余弦相似度 [-1.0, 1.0]
func CosineSimilarity(v1, v2 []float32) float32 {
	if len(v1) != len(v2) || len(v1) == 0 {
		return 0.0
	}
	var dotProduct, normA, normB float32
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
	}
	for i := 0; i < len(v1); i++ {
		normA += v1[i] * v1[i]
		normB += v2[i] * v2[i]
	}
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// ============================================================================
// 2. 向量数据库索引与近似近邻检索 (ANN Vector Index)
// ============================================================================

type Document struct {
	ID        int64
	Content   string
	Embedding []float32 // 稠密向量 (Dense Vector)
	Keywords  []string  // 稀疏关键词 (Sparse Token)
}

type SearchResult struct {
	Doc   *Document
	Score float32
}

type VectorStore struct {
	docs []*Document
	mu   sync.RWMutex
}

func NewVectorStore() *VectorStore {
	return &VectorStore{docs: make([]*Document, 0)}
}

func (vs *VectorStore) Insert(doc *Document) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.docs = append(vs.docs, doc)
}

// VectorSearch 稠密向量 Top-K 语义相似度召回
func (vs *VectorStore) VectorSearch(queryVector []float32, topK int) []SearchResult {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	results := make([]SearchResult, 0, len(vs.docs))
	for _, doc := range vs.docs {
		sim := CosineSimilarity(queryVector, doc.Embedding)
		results = append(results, SearchResult{Doc: doc, Score: sim})
	}

	// 降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > topK {
		results = results[:topK]
	}
	return results
}

// ============================================================================
// 3. 稀疏检索与混合重排 (Hybrid Search & RRF Reciprocal Rank Fusion)
// ============================================================================

// KeywordSearch 基于 BM25/词频的字面稀疏召回
func (vs *VectorStore) KeywordSearch(queryTerms []string, topK int) []SearchResult {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	results := make([]SearchResult, 0)
	for _, doc := range vs.docs {
		var matchCount float32
		for _, term := range queryTerms {
			for _, kw := range doc.Keywords {
				if strings.EqualFold(kw, term) {
					matchCount += 1.0
				}
			}
		}
		if matchCount > 0 {
			results = append(results, SearchResult{Doc: doc, Score: matchCount})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > topK {
		results = results[:topK]
	}
	return results
}

// HybridSearchWithRRF 混合召回 + RRF (倒数排名融合) 重排
// RRF 公式: Score = 1 / (k + rank)
func (vs *VectorStore) HybridSearchWithRRF(queryVector []float32, queryTerms []string, topK int, rrfK float32) []SearchResult {
	vectorRankings := vs.VectorSearch(queryVector, topK*2)
	keywordRankings := vs.KeywordSearch(queryTerms, topK*2)

	scoreMap := make(map[int64]float32)
	docMap := make(map[int64]*Document)

	for rank, r := range vectorRankings {
		docMap[r.Doc.ID] = r.Doc
		scoreMap[r.Doc.ID] += 1.0 / (rrfK + float32(rank+1))
	}

	for rank, r := range keywordRankings {
		docMap[r.Doc.ID] = r.Doc
		scoreMap[r.Doc.ID] += 1.0 / (rrfK + float32(rank+1))
	}

	finalResults := make([]SearchResult, 0, len(docMap))
	for id, score := range scoreMap {
		finalResults = append(finalResults, SearchResult{Doc: docMap[id], Score: score})
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Score > finalResults[j].Score
	})

	if len(finalResults) > topK {
		finalResults = finalResults[:topK]
	}
	return finalResults
}

// ============================================================================
// 4. RAG Pipeline 知识组装与上下文构建
// ============================================================================

type RAGPipeline struct {
	store *VectorStore
}

func NewRAGPipeline(store *VectorStore) *RAGPipeline {
	return &RAGPipeline{store: store}
}

// BuildAugmentedPrompt 召回知识库文档并组装最终 LLM Prompt
func (p *RAGPipeline) BuildAugmentedPrompt(ctx context.Context, userQuery string, queryVec []float32, queryTerms []string) (string, error) {
	hits := p.store.HybridSearchWithRRF(queryVec, queryTerms, 2, 60.0)

	var sb strings.Builder
	sb.WriteString("你是一个专业的技术知识助手。请参考以下参考背景信息回答用户的问题：\n\n【参考资料】：\n")
	for i, hit := range hits {
		sb.WriteString(fmt.Sprintf("[%d] (匹配得分: %.4f): %s\n", i+1, hit.Score, hit.Doc.Content))
	}
	sb.WriteString(fmt.Sprintf("\n【用户提问】：%s\n【回答】：", userQuery))
	return sb.String(), nil
}

func main() {
	fmt.Println("=== 阶段四 专题02：向量数据库与 RAG 混合召回重排链路 ===")

	store := NewVectorStore()
	// 插入知识库测试数据 (4维玩具向量)
	store.Insert(&Document{
		ID:        1,
		Content:   "Go 语言 GMP 调度器通过 Work-stealing 和 Sysmon 保证协程调度高吞吐",
		Embedding: []float32{0.9, 0.1, 0.2, 0.0},
		Keywords:  []string{"GMP", "调度器", "协程"},
	})
	store.Insert(&Document{
		ID:        2,
		Content:   "Redis 采用单线程 Reactor 模型配合跳表实现高性能分布式缓存与锁",
		Embedding: []float32{0.1, 0.9, 0.1, 0.3},
		Keywords:  []string{"Redis", "缓存", "分布式锁"},
	})
	store.Insert(&Document{
		ID:        3,
		Content:   "Go 运行时 GC 采用三色标记法结合混合写屏障，大幅消灭 STW 停顿",
		Embedding: []float32{0.8, 0.2, 0.3, 0.1},
		Keywords:  []string{"GC", "三色标记", "STW"},
	})

	pipeline := NewRAGPipeline(store)

	queryText := "请解释 Go 调度器的工作原理"
	queryVec := []float32{0.88, 0.12, 0.22, 0.0}
	queryTerms := []string{"GMP", "调度器"}

	prompt, _ := pipeline.BuildAugmentedPrompt(context.Background(), queryText, queryVec, queryTerms)
	fmt.Println(prompt)
}
