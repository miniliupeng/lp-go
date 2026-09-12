package main

import (
	"context"
	"strings"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	v1 := []float32{1.0, 0.0, 0.0}
	v2 := []float32{1.0, 0.0, 0.0}
	sim := CosineSimilarity(v1, v2)
	if sim < 0.999 {
		t.Fatalf("expected 1.0, got %f", sim)
	}

	v3 := []float32{0.0, 1.0, 0.0}
	simOrthogonal := CosineSimilarity(v1, v3)
	if simOrthogonal != 0.0 {
		t.Fatalf("expected 0.0 for orthogonal vectors, got %f", simOrthogonal)
	}
}

func TestRAGPipelineRetrieval(t *testing.T) {
	store := NewVectorStore()
	store.Insert(&Document{
		ID:        1,
		Content:   "Go GMP 协程调度实战",
		Embedding: []float32{0.9, 0.1},
		Keywords:  []string{"GMP", "调度"},
	})
	store.Insert(&Document{
		ID:        2,
		Content:   "Kafka 消息队列高可用",
		Embedding: []float32{0.1, 0.9},
		Keywords:  []string{"Kafka", "MQ"},
	})

	pipeline := NewRAGPipeline(store)
	prompt, err := pipeline.BuildAugmentedPrompt(context.Background(), "怎么学 GMP", []float32{0.85, 0.15}, []string{"GMP"})
	if err != nil {
		t.Fatalf("failed to build prompt: %v", err)
	}

	if !strings.Contains(prompt, "Go GMP 协程调度实战") {
		t.Fatalf("expected document 1 in prompt, got: %s", prompt)
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	// 模拟主流 1536 维 (OpenAI text-embedding-ada-002) 向量
	v1 := make([]float32, 1536)
	v2 := make([]float32, 1536)
	for i := 0; i < 1536; i++ {
		v1[i] = float32(i) * 0.001
		v2[i] = float32(1536-i) * 0.001
	}

	b.ResetTimer()
	for b.Loop() {
		_ = CosineSimilarity(v1, v2)
	}
}

func BenchmarkHybridSearchWithRRF(b *testing.B) {
	store := NewVectorStore()
	for i := 0; i < 100; i++ {
		store.Insert(&Document{
			ID:        int64(i),
			Content:   "Sample Vector Doc",
			Embedding: []float32{float32(i) * 0.01, 0.5, 0.2, 0.1},
			Keywords:  []string{"Go", "RAG", "AI"},
		})
	}

	queryVec := []float32{0.5, 0.5, 0.2, 0.1}
	queryTerms := []string{"Go"}

	b.ResetTimer()
	for b.Loop() {
		_ = store.HybridSearchWithRRF(queryVec, queryTerms, 5, 60.0)
	}
}
