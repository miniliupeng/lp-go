package upstream

import (
	"context"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
	"lp-go/07-capstone-project/ai-gateway-service/internal/streaming"
)

// LLMClient 上游大模型服务调用统一抽象接口
type LLMClient interface {
	// Complete 同步非流式推理
	Complete(ctx context.Context, req *domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error)
	// StreamComplete 异步流式推理，通过有界管道持续推送 Token
	StreamComplete(ctx context.Context, req *domain.ChatCompletionRequest, pipe *streaming.BackpressureChannel) error
}
