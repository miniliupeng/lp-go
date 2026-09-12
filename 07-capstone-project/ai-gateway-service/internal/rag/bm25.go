package rag

import (
	"sort"
	"strings"
	"sync"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
)

// BM25Index 稀疏关键词倒排索引 (字面精确匹配)
type BM25Index struct {
	docs []*domain.Document
	mu   sync.RWMutex
}

func NewBM25Index() *BM25Index {
	return &BM25Index{docs: make([]*domain.Document, 0)}
}

func (bi *BM25Index) Add(doc *domain.Document) {
	bi.mu.Lock()
	defer bi.mu.Unlock()
	bi.docs = append(bi.docs, doc)
}

func (bi *BM25Index) SearchTopK(keywords []string, topK int) []domain.SearchHit {
	bi.mu.RLock()
	defer bi.mu.RUnlock()

	hits := make([]domain.SearchHit, 0)
	for _, doc := range bi.docs {
		var score float32
		for _, kw := range keywords {
			for _, dkw := range doc.Keywords {
				if strings.EqualFold(kw, dkw) {
					score += 1.0
				}
			}
		}
		if score > 0 {
			hits = append(hits, domain.SearchHit{Doc: doc, Score: score})
		}
	}

	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Score > hits[j].Score
	})

	if len(hits) > topK {
		hits = hits[:topK]
	}
	return hits
}
