package main

import (
	"testing"
)

func TestInspectEnvironment(t *testing.T) {
	info := InspectEnvironment()

	if !info.IsHealthy {
		t.Fatalf("预期环境健康状态为 true，实际为 false")
	}

	if info.GoVersion == "" {
		t.Errorf("GoVersion 不应为空")
	}

	if info.GoOS == "" {
		t.Errorf("GoOS 操作系统标识不应为空")
	}

	if info.GoArch == "" {
		t.Errorf("GoArch 架构标识不应为空")
	}

	if info.NumCPU <= 0 {
		t.Errorf("CPU 逻辑核心数必须大于 0，实际为: %d", info.NumCPU)
	}
}
