package main

import (
	"fmt"
)

// CalculateFibonacci 计算斐波那契数列第 n 项，适合作为单步调试断点演示
func CalculateFibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		// 可以在此处打断点，单步观察 prev 与 curr 的值变化
		prev, curr = curr, prev+curr
	}
	return curr
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("       🛠️ 现代开发工具链体验与算法调试演示        ")
	fmt.Println("==================================================")

	target := 8
	result := CalculateFibonacci(target)

	fmt.Printf("计算斐波那契数列第 %d 项结果为: %d\n", target, result)
	fmt.Println("💡 [调试提示]: 在 CalculateFibonacci 内部打断点并按 F5，即可单步观察内存变量！")
	fmt.Println("==================================================")
}
