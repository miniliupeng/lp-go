package main

import (
	"fmt"
	"os"
	"runtime"
)

// EnvironmentInfo 封装当前 Go 开发环境的核心参数
type EnvironmentInfo struct {
	GoVersion   string
	GoOS        string
	GoArch      string
	Compiler    string
	NumCPU      int
	GOPROXY     string
	GOROOT      string
	IsHealthy   bool
}

// InspectEnvironment 收集并诊断当前 Go 环境健康状态
func InspectEnvironment() EnvironmentInfo {
	goproxy := os.Getenv("GOPROXY")
	if goproxy == "" {
		goproxy = "未显式设置系统环境变量 (使用 go env 内部默认值)"
	}

	goroot := runtime.GOROOT()
	if goroot == "" {
		goroot = "系统默认 SDK 路径"
	}

	info := EnvironmentInfo{
		GoVersion: runtime.Version(),
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		Compiler:  runtime.Compiler,
		NumCPU:    runtime.NumCPU(),
		GOPROXY:   goproxy,
		GOROOT:    goroot,
		IsHealthy: runtime.NumCPU() > 0 && runtime.Version() != "",
	}

	return info
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("       🚀 Go 语言现代开发环境诊断与健康自检器       ")
	fmt.Println("==================================================")

	info := InspectEnvironment()

	fmt.Printf("1. 当前 Go 运行时版本 : %s\n", info.GoVersion)
	fmt.Printf("2. 宿主操作系统 (GOOS): %s\n", info.GoOS)
	fmt.Printf("3. 目标硬件架构(GOARCH): %s\n", info.GoArch)
	fmt.Printf("4. 默认编译器标识     : %s\n", info.Compiler)
	fmt.Printf("5. 宿主可用 CPU 核心数: %d 逻辑核\n", info.NumCPU)
	fmt.Printf("6. GOROOT SDK 根目录  : %s\n", info.GOROOT)
	fmt.Printf("7. GOPROXY 模块代理   : %s\n", info.GOPROXY)

	fmt.Println("--------------------------------------------------")
	if info.IsHealthy {
		fmt.Println("✅ [诊断结论]: Go 开发环境状态正常，SDK 可正常编译执行！")
	} else {
		fmt.Println("❌ [诊断结论]: 环境存在异常，请检查 Go 安装完整性。")
	}
	fmt.Println("==================================================")
}
