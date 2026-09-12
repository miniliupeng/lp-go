package main

import (
	"strings"
	"testing"
)

func TestGenerateGreeting(t *testing.T) {
	// 测试默认参数
	defaultGreeting := GenerateGreeting([]string{"program"})
	if !strings.Contains(defaultGreeting, "Gopher") {
		t.Errorf("未传参数时应默认问候 Gopher，实际输出: %s", defaultGreeting)
	}

	// 测试自定义参数
	customGreeting := GenerateGreeting([]string{"program", "Alice"})
	if !strings.Contains(customGreeting, "Alice") {
		t.Errorf("传入 Alice 时应问候 Alice，实际输出: %s", customGreeting)
	}
}
