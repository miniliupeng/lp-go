package upstream

import (
	"context"
	"fmt"
	"strings"
	"time"

	"lp-go/07-capstone-project/ai-gateway-service/internal/domain"
	"lp-go/07-capstone-project/ai-gateway-service/internal/streaming"
)

// MockEngine 仿真大模型推理引擎，支持模拟生成延迟、断开取消及故障注入
type MockEngine struct {
	tokenDelay time.Duration
	failAll    bool
}

func NewMockEngine(tokenDelay time.Duration) *MockEngine {
	if tokenDelay <= 0 {
		tokenDelay = 10 * time.Millisecond
	}
	return &MockEngine{tokenDelay: tokenDelay}
}

func (m *MockEngine) SetFail(fail bool) {
	m.failAll = fail
}

func (m *MockEngine) Complete(ctx context.Context, req *domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	if m.failAll {
		return nil, fmt.Errorf("upstream mock engine failure: service unavailable")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(m.tokenDelay * 3):
	}

	lastPrompt := ""
	if len(req.Messages) > 0 {
		lastPrompt = req.Messages[len(req.Messages)-1].Content
	}

	return &domain.ChatCompletionResponse{
		ID:      "chatcmpl-mock-single-001",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []domain.ChatChoice{
			{
				Index: 0,
				Message: domain.ChatMessage{
					Role:    domain.RoleAssistant,
					Content: fmt.Sprintf("[Mock Response for: %s]", lastPrompt),
				},
				FinishReason: "stop",
			},
		},
	}, nil
}

func (m *MockEngine) StreamComplete(ctx context.Context, req *domain.ChatCompletionRequest, pipe *streaming.BackpressureChannel) error {
	defer pipe.Close()

	if m.failAll {
		return fmt.Errorf("upstream stream error: connection refused")
	}

	lastPrompt := ""
	if len(req.Messages) > 0 {
		lastPrompt = req.Messages[len(req.Messages)-1].Content
	}

	tokens := []string{"【AI网关】", "收到", "您的", "提问", "：", lastPrompt, "。", "正在", "以高", "性能", "流式", "输出", "完成", "响应", "。"}

	for i, tok := range tokens {
		select {
		case <-ctx.Done():
			// 核心：客户端主动断连后，即刻退出协程，终止后续 Token 模拟与算力消耗
			return ctx.Err()
		case <-time.After(m.tokenDelay):
			var finishReason *string
			if i == len(tokens)-1 {
				s := "stop"
				finishReason = &s
			}

			chunk := &domain.ChatCompletionChunk{
				ID:      "chatcmpl-stream-001",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   req.Model,
				Choices: []domain.StreamChoice{
					{
						Index: 0,
						Delta: domain.StreamDelta{
							Role:    domain.RoleAssistant,
							Content: tok,
						},
						FinishReason: finishReason,
					},
				},
			}

			// 写入反压管道，超时 100ms 则触发反压保护
			if err := pipe.Push(ctx, chunk, 100*time.Millisecond); err != nil {
				return err
			}
		}
	}
	_ = strings.Clone("")
	return nil
}
