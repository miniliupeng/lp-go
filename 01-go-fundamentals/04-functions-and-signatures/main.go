package main

import (
	"errors"
	"fmt"
)

// SafeDivide 演示多返回值规范
func SafeDivide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零 (division by zero)")
	}
	return a / b, nil
}

// CalculateStatistics 演示命名返回值与变参函数
func CalculateStatistics(scores ...float64) (min float64, max float64, avg float64) {
	if len(scores) == 0 {
		return 0, 0, 0
	}

	min = scores[0]
	max = scores[0]
	total := 0.0

	for _, s := range scores {
		if s < min {
			min = s
		}
		if s > max {
			max = s
		}
		total += s
	}

	avg = total / float64(len(scores))
	return min, max, avg
}

// MakeSequenceGenerator 演示闭包状态捕获机制
func MakeSequenceGenerator(start int, step int) func() int {
	current := start
	return func() int {
		val := current
		current += step
		return val
	}
}

// modifyValueAttempt 尝试修改传入的值（验证值拷贝）
func modifyValueAttempt(x int) {
	x = 999
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("       📐 专题 04：函数签名、多返回值与闭包实操     ")
	fmt.Println("==================================================")

	// 1. 多返回值测试
	val, err := SafeDivide(10, 2)
	fmt.Printf("1. 10 / 2 = %.2f, 错误状态: %v\n", val, err)
	_, errZero := SafeDivide(10, 0)
	fmt.Printf("   10 / 0 触发错误返回: %v\n", errZero)

	// 2. 变参函数与命名返回值
	min, max, avg := CalculateStatistics(85.5, 92.0, 78.0, 99.5, 60.0)
	fmt.Printf("\n2. 统计变参计算结果: 最低分=%.1f, 最高分=%.1f, 平均分=%.1f\n", min, max, avg)

	// 3. 闭包序列生成器
	nextEven := MakeSequenceGenerator(0, 2)
	fmt.Printf("\n3. 闭包捕获外部变量连续步进: %d -> %d -> %d -> %d\n",
		nextEven(), nextEven(), nextEven(), nextEven())

	// 4. 参数值拷贝证明
	original := 42
	modifyValueAttempt(original)
	fmt.Printf("\n4. 传参值拷贝验证: 传入原变量 42，函数内修改后外部变量依然为: %d\n", original)

	fmt.Println("==================================================")
}
