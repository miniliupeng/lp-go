package rag

import (
	"sort"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
)

// RRFMerger 倒数排名融合算法 (Reciprocal Rank Fusion) 重排器
// 公式: Score(d) = \sum_{r \in Rankings} \frac{1}{k + rank(d)}
type RRFMerger struct {
	k float32
}

func NewRRFMerger(k float32) *RRFMerger {
	if k <= 0 {
		k = 60.0
	}
	return &RRFMerger{k: k}
}

func (m *RRFMerger) Merge(denseRankings, sparseRankings []domain.SearchHit, topK int) []domain.SearchHit {
	scoreMap := make(map[int64]float32)
	docMap := make(map[int64]*domain.Document)

	for rank, hit := range denseRankings {
		docMap[hit.Doc.ID] = hit.Doc
		scoreMap[hit.Doc.ID] += 1.0 / (m.k + float32(rank+1))
	}

	for rank, hit := range sparseRankings {
		docMap[hit.Doc.ID] = hit.Doc
		scoreMap[hit.Doc.ID] += 1.0 / (m.k + float32(rank+1))
	}

	merged := make([]domain.SearchHit, 0, len(docMap))
	for id, score := range scoreMap {
		merged = append(merged, domain.SearchHit{Doc: docMap[id], Score: score})
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Score > merged[j].Score
	})

	if len(merged) > topK {
		merged = merged[:topK]
	}
	return merged
}
