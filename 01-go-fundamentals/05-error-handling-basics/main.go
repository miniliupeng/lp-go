package main

import (
	"errors"
	"fmt"
)

// ProcessTransaction 演示显式错误契约
func ProcessTransaction(amount int, balance int) (newBalance int, err error) {
	if amount <= 0 {
		return balance, errors.New("交易扣款金额必须大于 0")
	}
	if amount > balance {
		return balance, fmt.Errorf("账户余额不足: 当前余额 %d，请求扣款 %d", balance, amount)
	}
	return balance - amount, nil
}

// SafeRunWorker 演示利用 defer 和 recover 构建兜底防崩溃护盾
func SafeRunWorker(shouldPanic bool) (intercepted bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("   🛡️ [Recovery 中间件生效]: 成功捕获 Panic: %v\n", r)
			intercepted = true
		}
	}()

	if shouldPanic {
		panic("模拟发生未知空指针或越界 Panic 崩溃！")
	}
	return false
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("    🛡️ 专题 05：显式错误处理与 defer / recover    ")
	fmt.Println("==================================================")

	// 1. 显式业务错误返回验证
	balance := 100
	fmt.Printf("初始账户余额: %d 元\n", balance)

	newBal, err := ProcessTransaction(30, balance)
	if err != nil {
		fmt.Printf("扣款失败: %v\n", err)
	} else {
		fmt.Printf("1. 成功扣款 30 元，当前最新余额: %d 元\n", newBal)
		balance = newBal
	}

	_, errOver := ProcessTransaction(200, balance)
	if errOver != nil {
		fmt.Printf("2. 超额扣款预期触发拦截: %v\n", errOver)
	}

	// 2. defer LIFO 时序验证
	fmt.Println("\n3. 验证多个 defer 注册时的 LIFO（后进先出）时序:")
	func() {
		defer fmt.Println("   第 1 个 defer: [资源 A 释放]")
		defer fmt.Println("   第 2 个 defer: [资源 B 释放]")
		defer fmt.Println("   第 3 个 defer: [资源 C 释放]")
		fmt.Println("   -> 函数核心工作完成，即将退出...")
	}()

	// 3. recover 异常兜底防护
	fmt.Println("\n4. 验证 panic 与 recover 异常兜底:")
	caught := SafeRunWorker(true)
	fmt.Printf("   Worker 是否被成功兜底保护且进程未崩溃: %v\n", caught)

	fmt.Println("==================================================")
}
