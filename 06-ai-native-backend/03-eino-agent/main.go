package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 1. Eino 风格核心类型体系 (Message / Tool / Invocation)
// ============================================================================

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type ToolCall struct {
	ID        string `json:"id"`
	Function  string `json:"function"`
	Arguments string `json:"arguments"`
}

type Message struct {
	Role       Role        `json:"role"`
	Content    string      `json:"content"`
	ToolCalls  []*ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

// Tool 工具接口规范
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, argsJSON string) (string, error)
}

// ============================================================================
// 2. 真实工具实现 (WeatherTool / CalculatorTool)
// ============================================================================

type WeatherTool struct{}

func (w *WeatherTool) Name() string { return "get_weather" }
func (w *WeatherTool) Description() string {
	return "查询指定城市的天气状况，参数格式: {\"city\": \"Beijing\"}"
}

func (w *WeatherTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var payload struct {
		City string `json:"city"`
	}
	if err := json.Unmarshal([]byte(argsJSON), &payload); err != nil {
		return "", err
	}
	if payload.City == "" {
		return "", errors.New("city parameter is required")
	}
	return fmt.Sprintf("{\"city\": \"%s\", \"weather\": \"晴空万里\", \"temperature\": \"24°C\"}", payload.City), nil
}

type CalculatorTool struct{}

func (c *CalculatorTool) Name() string { return "calculator" }
func (c *CalculatorTool) Description() string {
	return "基础数学四则计算器，参数格式: {\"expression\": \"100 * 25\"}"
}

func (c *CalculatorTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var payload struct {
		Expression string `json:"expression"`
	}
	if err := json.Unmarshal([]byte(argsJSON), &payload); err != nil {
		return "", err
	}
	// 演示简易计算逻辑
	return "{\"expression\": \"" + payload.Expression + "\", \"result\": 2500}", nil
}

// ============================================================================
// 3. Eino 原生 Graph 节点与 ReAct (Reasoning + Acting) Agent 状态循环
// ============================================================================

type MockLLMClient struct{}

// Generate 模拟大模型的 Tool Calling 意图识别与回答决策
func (m *MockLLMClient) Generate(ctx context.Context, history []*Message) (*Message, error) {
	if len(history) == 0 {
		return nil, errors.New("empty history")
	}

	lastMsg := history[len(history)-1]

	// 1. 若上一条是工具执行结果，则大模型基于上下文总结输出最终回答
	if lastMsg.Role == RoleTool {
		return &Message{
			Role:    RoleAssistant,
			Content: fmt.Sprintf("根据工具返回的最新数据，为您提供最终结论：%s", lastMsg.Content),
		}, nil
	}

	// 2. 模拟根据用户输入判定需要调用天气工具
	if strings.Contains(lastMsg.Content, "天气") {
		return &Message{
			Role: RoleAssistant,
			ToolCalls: []*ToolCall{
				{
					ID:        "call_weather_9527",
					Function:  "get_weather",
					Arguments: `{"city": "Beijing"}`,
				},
			},
		}, nil
	}

	// 3. 普通文本直出
	return &Message{
		Role:    RoleAssistant,
		Content: "您好！我是基于字节 Eino 架构构建的通用智能助理，请问有什么可以帮您？",
	}, nil
}

// EinoAgent 基于 Graph 状态图运行的 ReAct 智能体
type EinoAgent struct {
	llm   *MockLLMClient
	tools map[string]Tool
	mu    sync.RWMutex
}

func NewEinoAgent() *EinoAgent {
	agent := &EinoAgent{
		llm:   &MockLLMClient{},
		tools: make(map[string]Tool),
	}
	// 注册开箱即用工具
	agent.RegisterTool(&WeatherTool{})
	agent.RegisterTool(&CalculatorTool{})
	return agent
}

func (a *EinoAgent) RegisterTool(tool Tool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tools[tool.Name()] = tool
}

// Run 状态图循环（Loop MaxIterations 保护）
func (a *EinoAgent) Run(ctx context.Context, userPrompt string, maxIterations int) ([]*Message, error) {
	messages := []*Message{
		{Role: RoleUser, Content: userPrompt},
	}

	for i := 0; i < maxIterations; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 1. Agent 思考节点 (LLM Node)
		reply, err := a.llm.Generate(ctx, messages)
		if err != nil {
			return nil, err
		}
		messages = append(messages, reply)

		// 2. 判定是否有 Tool Calling 需求
		if len(reply.ToolCalls) == 0 {
			// 无需工具调用，已得到最终答复，终止循环
			return messages, nil
		}

		// 3. 工具执行节点 (Tool Execution Node)
		for _, toolCall := range reply.ToolCalls {
			a.mu.RLock()
			tool, ok := a.tools[toolCall.Function]
			a.mu.RUnlock()

			var toolResult string
			if !ok {
				toolResult = fmt.Sprintf("Error: tool %s not found", toolCall.Function)
			} else {
				res, err := tool.Execute(ctx, toolCall.Arguments)
				if err != nil {
					toolResult = fmt.Sprintf("Error: %v", err)
				} else {
					toolResult = res
				}
			}

			// 将工具执行结果作为 RoleTool 消息推入上下文
			messages = append(messages, &Message{
				Role:       RoleTool,
				ToolCallID: toolCall.ID,
				Content:    toolResult,
			})
		}
	}

	return nil, errors.New("agent reached max iterations limit")
}

func main() {
	fmt.Println("=== 阶段四 专题03：字节 Eino 框架落地（Agent 工具调用与 Graph 编排） ===")

	agent := NewEinoAgent()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userQuery := "请帮我查一下北京今天的天气如何？"
	fmt.Printf("[用户提问]: %s\n\n", userQuery)

	history, err := agent.Run(ctx, userQuery, 5)
	if err != nil {
		fmt.Printf("Agent 执行失败: %v\n", err)
		return
	}

	for _, msg := range history {
		if len(msg.ToolCalls) > 0 {
			fmt.Printf("👉 [%s (Tool Calling 发起)]: 触发工具 %s, 入参 %s\n",
				msg.Role, msg.ToolCalls[0].Function, msg.ToolCalls[0].Arguments)
		} else if msg.Role == RoleTool {
			fmt.Printf("🔧 [%s (工具执行反馈)]: %s\n", msg.Role, msg.Content)
		} else {
			fmt.Printf("💬 [%s]: %s\n", msg.Role, msg.Content)
		}
	}
}
