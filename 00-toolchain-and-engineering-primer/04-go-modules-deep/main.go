package main

import (
	"crypto/rand"
	"fmt"
)

// ModuleInfo 模拟模块元数据结构体（大写首字母导出）
type ModuleInfo struct {
	Name    string
	Version string
	Active  bool
}

// GeneratePseudoUUID 模拟纯 Go 生成轻量级唯一序列标识符，展示标准库包依赖封装
func GeneratePseudoUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("       📦 专题 04：Go Modules 依赖与工程封装演示      ")
	fmt.Println("==================================================")

	mod := ModuleInfo{
		Name:    "lp-go",
		Version: "v1.0.0-2026",
		Active:  true,
	}

	fmt.Printf("当前工程根模块名 : %s\n", mod.Name)
	fmt.Printf("当前工程发布版本 : %s\n", mod.Version)

	requestID := GeneratePseudoUUID()
	fmt.Printf("模拟生成的业务 Request-ID : %s\n", requestID)

	fmt.Println("\n💡 [依赖治理军规]:")
	fmt.Println("  1. 引入外部包后，务必执行 `go mod tidy` 自动校准依赖！")
	fmt.Println("  2. `go.sum` 必须提交到 Git 版本控制中，确保团队构建一致性。")
	fmt.Println("==================================================")
}
