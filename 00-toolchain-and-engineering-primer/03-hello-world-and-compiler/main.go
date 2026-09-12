package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateGreeting 根据输入参数生成规范问候语
func GenerateGreeting(args []string) string {
	target := "Gopher"
	if len(args) > 1 {
		target = args[1]
	}
	return fmt.Sprintf("Hello, %s! 欢迎开启 Go 全栈工程进阶之旅！", target)
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("      🌟 专题 03：第一个 Go 程序与 CLI 参数解析     ")
	fmt.Println("==================================================")

	// 获取当前可执行程序路径与命令行参数
	execPath := os.Args[0]
	execName := filepath.Base(execPath)

	fmt.Printf("当前执行进程名称 : %s\n", execName)
	fmt.Printf("接收到命令行参数数量: %d\n", len(os.Args))

	// 生成问候语
	greeting := GenerateGreeting(os.Args)
	fmt.Println("\n💬 [程序回显]:")
	fmt.Println("  ", greeting)

	fmt.Println("\n💡 [实操提示]: 尝试使用不同参数运行:")
	fmt.Println("   go run ./00-toolchain-and-engineering-primer/03-hello-world-and-compiler/main.go \"未来架构师\"")
	fmt.Println("==================================================")
}
