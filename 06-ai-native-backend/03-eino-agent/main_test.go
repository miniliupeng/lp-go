package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestEinoAgentWeatherToolCalling(t *testing.T) {
	agent := NewEinoAgent()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	history, err := agent.Run(ctx, "北京天气怎么样？", 5)
	if err != nil {
		t.Fatalf("agent run failed: %v", err)
	}

	// 验证消息流转是否包含 Tool 调用与 Tool 结果
	hasToolCall := false
	hasToolResult := false
	hasFinalAssistant := false

	for _, msg := range history {
		if len(msg.ToolCalls) > 0 && msg.ToolCalls[0].Function == "get_weather" {
			hasToolCall = true
		}
		if msg.Role == RoleTool && strings.Contains(msg.Content, "晴空万里") {
			hasToolResult = true
		}
		if msg.Role == RoleAssistant && strings.Contains(msg.Content, "最终结论") {
			hasFinalAssistant = true
		}
	}

	if !hasToolCall {
		t.Fatal("expected tool call 'get_weather'")
	}
	if !hasToolResult {
		t.Fatal("expected tool result with weather info")
	}
	if !hasFinalAssistant {
		t.Fatal("expected final assistant conclusion")
	}
}

func TestEinoAgentDirectText(t *testing.T) {
	agent := NewEinoAgent()

	ctx := context.Background()
	history, err := agent.Run(ctx, "你好，打个招呼吧", 5)
	if err != nil {
		t.Fatalf("agent run failed: %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 messages (user + assistant), got %d", len(history))
	}
}

func BenchmarkWeatherToolExecution(b *testing.B) {
	tool := &WeatherTool{}
	ctx := context.Background()
	args := `{"city": "Beijing"}`

	b.ResetTimer()
	for b.Loop() {
		_, err := tool.Execute(ctx, args)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEinoAgentRun(b *testing.B) {
	agent := NewEinoAgent()
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, err := agent.Run(ctx, "你好", 3)
		if err != nil {
			b.Fatal(err)
		}
	}
}
