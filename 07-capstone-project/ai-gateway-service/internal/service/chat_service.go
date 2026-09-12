package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
	"lp-go/07-capstone-project/ai-gateway-service/internal/rag"
	"lp-go/07-capstone-project/ai-gateway-service/internal/streaming"
	"lp-go/07-capstone-project/ai-gateway-service/internal/upstream"
	"lp-go/07-capstone-project/ai-gateway-service/pkg/metrics"
)

// SingleFlight 极简并发请求合并器 (消除相同 Prompt 瞬时击穿)
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

// ChatService 网关核心业务协调层
type ChatService struct {
	upstream    upstream.LLMClient
	vectorIndex *rag.VectorIndex
	bm25Index   *rag.BM25Index
	rrfMerger   *rag.RRFMerger
	sfGroup     Group
	collector   *metrics.REDCollector
}

func NewChatService(up upstream.LLMClient, collector *metrics.REDCollector) *ChatService {
	svc := &ChatService{
		upstream:    up,
		vectorIndex: rag.NewVectorIndex(),
		bm25Index:   rag.NewBM25Index(),
		rrfMerger:   rag.NewRRFMerger(60.0),
		collector:   collector,
	}
	svc.seedRAGData()
	return svc
}

func (s *ChatService) seedRAGData() {
	// 预置企业私域知识库 (例如大厂内部 Go 网关架构规范)
	doc1 := &domain.Document{
		ID:        1,
		Title:     "网关高可用手册",
		Content:   "AI 网关生产规范要求必须开启有界通道网络反压，并在客户端断连时监听 ctx.Done() 取消上游计费。",
		Embedding: []float32{0.9, 0.1, 0.2, 0.05},
		Keywords:  []string{"反压", "网关", "ctx.Done()", "计费"},
	}
	doc2 := &domain.Document{
		ID:        2,
		Title:     "SingleFlight 防击穿实践",
		Content:   "在面对瞬时突发热点问答时，使用 SingleFlight 合并相同 Query，能大幅降低上游 GPU 推理卡压力。",
		Embedding: []float32{0.1, 0.85, 0.3, 0.2},
		Keywords:  []string{"SingleFlight", "击穿", "合并", "GPU"},
	}
	s.vectorIndex.Add(doc1)
	s.vectorIndex.Add(doc2)
	s.bm25Index.Add(doc1)
	s.bm25Index.Add(doc2)
}

// Complete 处理非流式请求 (集成 SingleFlight 防击穿与 RAG)
func (s *ChatService) Complete(ctx context.Context, req *domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	if req.EnableRAG {
		s.applyRAG(req)
	}

	cacheKey := req.Model
	if len(req.Messages) > 0 {
		cacheKey += ":" + req.Messages[len(req.Messages)-1].Content
	}

	res, err := s.sfGroup.Do(cacheKey, func() (any, error) {
		return s.upstream.Complete(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return res.(*domain.ChatCompletionResponse), nil
}

// Stream 处理 SSE 流式转发 (集成 RAG 上下文增强与 TTFT 统计)
func (s *ChatService) Stream(ctx context.Context, req *domain.ChatCompletionRequest, flusher *streaming.SSEFlusher) error {
	if req.EnableRAG {
		s.applyRAG(req)
	}

	pipe := streaming.NewBackpressureChannel(16)
	start := time.Now()
	firstTokenRecorded := false

	// 启动后台协程从上游拉取 Token
	go func() {
		_ = s.upstream.StreamComplete(ctx, req, pipe)
	}()

	// 消费管道输出并推送到前端
	for chunk := range pipe.Out() {
		if !firstTokenRecorded {
			s.collector.RecordTTFT(time.Since(start))
			firstTokenRecorded = true
		}

		if err := flusher.PushChunk(chunk); err != nil {
			return err
		}
	}

	return flusher.PushDone()
}

func (s *ChatService) applyRAG(req *domain.ChatCompletionRequest) {
	if len(req.Messages) == 0 {
		return
	}
	lastMsg := &req.Messages[len(req.Messages)-1]

	// 模拟查询向量与关键词提取
	queryVec := []float32{0.85, 0.15, 0.2, 0.1}
	queryKws := []string{"反压", "网关", "SingleFlight"}

	denseHits := s.vectorIndex.SearchTopK(queryVec, 2)
	sparseHits := s.bm25Index.SearchTopK(queryKws, 2)
	finalHits := s.rrfMerger.Merge(denseHits, sparseHits, 2)

	if len(finalHits) == 0 {
		return
	}

	var sb strings.Builder
	sb.WriteString("【系统私域参考资料】：\n")
	for i, hit := range finalHits {
		sb.WriteString(fmt.Sprintf("%d. 《%s》: %s\n", i+1, hit.Doc.Title, hit.Doc.Content))
	}
	sb.WriteString("\n【用户提问】：\n")
	sb.WriteString(lastMsg.Content)

	lastMsg.Content = sb.String()
}
