package rag

import (
	"math"
	"sort"
	"sync"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
)

// CosineSimilarity 计算高维浮点向量的余弦相似度
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}
	var dot, normA, normB float32
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

type VectorIndex struct {
	docs []*domain.Document
	mu   sync.RWMutex
}

func NewVectorIndex() *VectorIndex {
	return &VectorIndex{docs: make([]*domain.Document, 0)}
}

func (vi *VectorIndex) Add(doc *domain.Document) {
	vi.mu.Lock()
	defer vi.mu.Unlock()
	vi.docs = append(vi.docs, doc)
}

// SearchTopK 计算余弦相似度并按降序返回 Top-K
func (vi *VectorIndex) SearchTopK(queryVec []float32, topK int) []domain.SearchHit {
	vi.mu.RLock()
	defer vi.mu.RUnlock()

	hits := make([]domain.SearchHit, 0, len(vi.docs))
	for _, doc := range vi.docs {
		sim := CosineSimilarity(queryVec, doc.Embedding)
		hits = append(hits, domain.SearchHit{Doc: doc, Score: sim})
	}

	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Score > hits[j].Score
	})

	if len(hits) > topK {
		hits = hits[:topK]
	}
	return hits
}
