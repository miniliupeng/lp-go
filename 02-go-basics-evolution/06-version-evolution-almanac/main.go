package main

import (
	"errors"
	"fmt"
	"iter"
	"sync"
)

// 自定义哨兵错误
var ErrNotFound = errors.New("record not found")

// 自定义业务错误结构
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query [%s] failed: %v", e.Query, e.Err)
}

func (e *QueryError) Unwrap() error {
	return e.Err
}

func main() {
	fmt.Println("=== 1. Go 1.13+ 现代错误链 (errors.Is / As) ===")
	demoModernErrors()

	fmt.Println("\n=== 2. Go 1.18+ 现代泛型基础 (Generics) ===")
	demoGenericsBasics()

	fmt.Println("\n=== 3. Go 1.21+ 内置高效函数 (min/max/clear) ===")
	demoBuiltinFunctions()

	fmt.Println("\n=== 4. Go 1.22+ 循环变量独立作用域验证 ===")
	demoLoopVariableScope()

	fmt.Println("\n=== 5. Go 1.23+ 原生迭代器标准 (iter.Seq) ===")
	demoIterators()
}

// ------------------------------------------------------------
// 1. Go 1.13+ 错误处理演进
// ------------------------------------------------------------
func demoModernErrors() {
	baseErr := ErrNotFound
	wrappedErr := &QueryError{Query: "SELECT * FROM users", Err: fmt.Errorf("db error: %w", baseErr)}

	// errors.Is 沿链判定是否包含底层哨兵错误
	if errors.Is(wrappedErr, ErrNotFound) {
		fmt.Println("[errors.Is] 成功匹配到底层 ErrNotFound 哨兵错误")
	}

	// errors.As 沿链提取目标类型结构体
	var targetQueryErr *QueryError
	if errors.As(wrappedErr, &targetQueryErr) {
		fmt.Printf("[errors.As] 成功提取 QueryError 结构，查询语句: %s\n", targetQueryErr.Query)
	}
}

// ------------------------------------------------------------
// 2. Go 1.18+ 泛型基础
// ------------------------------------------------------------
// SumSlice 泛型切片求和函数（类型参数约束为 int 或 float64）
func SumSlice[T int | float64](items []T) T {
	var sum T
	for _, item := range items {
		sum += item
	}
	return sum
}

// Container 泛型结构体：承载任意类型的数据
type Container[T any] struct {
	Value T
}

func demoGenericsBasics() {
	intSum := SumSlice([]int{10, 20, 30})
	floatSum := SumSlice([]float64{1.5, 2.5, 4.0})
	fmt.Printf("泛型切片求和: intSum=%d, floatSum=%.1f\n", intSum, floatSum)

	strBox := Container[string]{Value: "现代化类型参数容器"}
	intBox := Container[int]{Value: 2026}
	fmt.Printf("泛型容器示例: strBox=%s, intBox=%d\n", strBox.Value, intBox.Value)
}

// ------------------------------------------------------------
// 3. Go 1.21+ 内置函数
// ------------------------------------------------------------
func demoBuiltinFunctions() {
	a, b := 42, 100
	fmt.Printf("min(%d, %d) = %d\n", a, b, min(a, b))
	fmt.Printf("max(%d, %d) = %d\n", a, b, max(a, b))

	m := map[string]int{"cpu": 8, "mem": 16}
	fmt.Printf("clear 前 map 大小: %d\n", len(m))
	clear(m) // 替代繁琐的遍历 delete，底层 memclr 极致优化
	fmt.Printf("clear 后 map 大小: %d\n", len(m))
}

// ------------------------------------------------------------
// 4. Go 1.22+ 循环变量独立作用域
// ------------------------------------------------------------
func demoLoopVariableScope() {
	values := []int{10, 20, 30}
	var wg sync.WaitGroup

	for _, v := range values {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Go 1.22 起每轮循环绑定全新变量实例，彻底消除旧版本闭包并发 Bug
			fmt.Printf("并发读取值: %d\n", v)
		}()
	}
	wg.Wait()
}

// ------------------------------------------------------------
// 5. Go 1.23+ 原生迭代器
// ------------------------------------------------------------
func demoIterators() {
	// 产生偶数序列的生成器
	evenNumbers := func(limit int) iter.Seq[int] {
		return func(yield func(int) bool) {
			for i := 0; i <= limit; i += 2 {
				if !yield(i) {
					return // 消费者 break 时提前退出
				}
			}
		}
	}

	fmt.Print("iter.Seq 产生的偶数: ")
	for val := range evenNumbers(10) {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}
