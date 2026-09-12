package domain

// RAG 领域实体契约

type Document struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"embedding"` // 稠密向量
	Keywords  []string  `json:"keywords"`  // 稀疏关键词
}

type SearchHit struct {
	Doc   *Document `json:"doc"`
	Score float32   `json:"score"`
}
