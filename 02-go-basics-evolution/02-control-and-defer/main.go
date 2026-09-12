package main

import (
	"errors"
	"fmt"
)

var ErrDatabase = errors.New("database connection lost")

type OpError struct {
	Op  string
	Err error
}

func (e *OpError) Error() string {
	return fmt.Sprintf("operation [%s] failed: %v", e.Op, e.Err)
}

func (e *OpError) Unwrap() error {
	return e.Err
}

func main() {
	fmt.Println("=== 1. 新手入门：流程控制与 1.22 循环语法糖 ===")
	demoControlFlow()

	fmt.Println("\n=== 2. 硬核底层：defer 汇编时序与命名返回值陷阱 ===")
	demoDeferReturnOrder()

	fmt.Println("\n=== 3. 硬核底层：for 循环中 defer 正确释放姿势 ===")
	demoLoopDeferPattern()

	fmt.Println("\n=== 4. 硬核底层：panic/recover 与现代 errors 链 ===")
	demoPanicAndModernErrors()
}

func demoControlFlow() {
	// 1. if 前置初始化
	if val := 42; val%2 == 0 {
		fmt.Printf("前置初始化: val = %d 是偶数\n", val)
	}

	// 2. Go 1.22+ 整数范围循环
	fmt.Print("1.22+ for i := range 4: ")
	for i := range 4 {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// 3. switch 自动 break
	status := 200
	switch status {
	case 200, 201:
		fmt.Println("switch 自动匹配成功，无需显式 break")
	default:
		fmt.Println("其他状态")
	}
}

// defer 时序验证 1：匿名返回值，返回 5
func deferF1() int {
	x := 5
	defer func() {
		x++ // 修改的是局部变量 x，不影响已被拷贝到返回值栈区的 5
	}()
	return x
}

// defer 时序验证 2：命名返回值，返回 6
func deferF2() (x int) {
	x = 5
	defer func() {
		x++ // x 就是返回值变量本身，直接被自增
	}()
	return x
}

// defer 时序验证 3：参数值传递，返回 5
func deferF3() (x int) {
	x = 5
	defer func(x int) {
		x++ // 参数同名变量遮蔽，修改的是入参局部变量
	}(x)
	return x
}

func demoDeferReturnOrder() {
	fmt.Printf("deferF1() 匿名返回值: %d (预期: 5)\n", deferF1())
	fmt.Printf("deferF2() 命名返回值: %d (预期: 6)\n", deferF2())
	fmt.Printf("deferF3() 参数值拷贝: %d (预期: 5)\n", deferF3())
}

// demoLoopDeferPattern 演示循环中安全 defer 模式
func demoLoopDeferPattern() {
	tasks := []string{"task_A", "task_B"}

	for _, task := range tasks {
		// 使用立即执行的匿名函数隔离 defer 作用域
		func() {
			fmt.Printf(">> 开始执行 %s...\n", task)
			defer fmt.Printf("<< %s 的资源在单次迭代结束时立即释放！\n", task)
		}()
	}
}

func demoPanicAndModernErrors() {
	// 1. 同协程 recover 保护
	safeCall := func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[recover 成功拦截] 发生 panic: %v\n", r)
			}
		}()
		panic("模拟发生未知空指针异常")
	}
	safeCall()

	// 2. Go 1.13+ errors 链式包装与拆解
	wrapped := &OpError{Op: "INSERT", Err: fmt.Errorf("network issue: %w", ErrDatabase)}
	if errors.Is(wrapped, ErrDatabase) {
		fmt.Println("[errors.Is] 成功在调用链深处识别到底层 ErrDatabase 哨兵错误！")
	}

	// 3. Go 1.20+ errors.Join
	err1 := errors.New("校验参数失败")
	err2 := errors.New("用户权限不足")
	allErrors := errors.Join(err1, err2)
	fmt.Printf("[errors.Join 聚合输出]:\n%v\n", allErrors)
}
