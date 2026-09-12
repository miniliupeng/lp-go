package main

import (
	"fmt"
)

// UpdateValueByPointer 演示通过指针直接修改调用方内存中的值
func UpdateValueByPointer(ptr *int, multiplier int) {
	if ptr != nil {
		*ptr = (*ptr) * multiplier
	}
}

// FilterEvenNumbers 演示切片 append 与遍历基础
func FilterEvenNumbers(src []int) []int {
	result := make([]int, 0, len(src))
	for _, n := range src {
		if n%2 == 0 {
			result = append(result, n)
		}
	}
	return result
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("   📦 专题 06：指针入门与切片/字典基础 CRUD 实操   ")
	fmt.Println("==================================================")

	// 1. 指针修改值实证
	val := 50
	fmt.Printf("1. 原值: %d (物理内存地址: %p)\n", val, &val)
	UpdateValueByPointer(&val, 3)
	fmt.Printf("   传入指针修改后，原变量值变为: %d\n", val)

	// 2. 数组与切片（含编译器自动推导数组长度 [...]int）
	arrInferred := [...]int{10, 20, 30} // 编译器自动推导固定长度为 [3]int
	fmt.Printf("\n2. 固定数组推导: %v (类型长度 len=%d)\n", arrInferred, len(arrInferred))

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	evens := FilterEvenNumbers(numbers)
	fmt.Printf("   切片过滤演示: 原切片=%v (cap=%d), 偶数切片=%v\n", numbers, cap(numbers), evens)

	// 3. 字典 Map 基础增删改查与遍历随机性
	userRegistry := make(map[string]string)
	userRegistry["u101"] = "张三"
	userRegistry["u102"] = "李四"
	userRegistry["u103"] = "王五"
	fmt.Printf("\n3. Map 写入后状态: %v\n", userRegistry)

	// comma-ok 安全读取
	if name, ok := userRegistry["u101"]; ok {
		fmt.Printf("   安全读取用户 u101: %s (ok=%v)\n", name, ok)
	}

	// 演示遍历 Map 键值对（注意：Go 运行时随机 hash seed 导致每次遍历顺序均不固定）
	fmt.Print("   Map range 遍历顺序: ")
	for k, v := range userRegistry {
		fmt.Printf("[%s:%s] ", k, v)
	}
	fmt.Println()

	// 4. Map comma-ok 读取与安全删除 (Delete)
	delete(userRegistry, "u102")
	fmt.Printf("\n4. 删除 u102 后的 Map: %v (元素总数 len=%d)\n", userRegistry, len(userRegistry))

	fmt.Println("==================================================")
}
