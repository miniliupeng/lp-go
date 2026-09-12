package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("=== 1. 新手入门：切片与 Map 常用操作 ===")
	demoBasicSliceAndMap()

	fmt.Println("\n=== 2. 硬核底层：三索引切片与内存隔离 ===")
	demoFullSliceExpression()

	fmt.Println("\n=== 3. 硬核底层：切片扩容容量平滑演变 ===")
	demoSliceGrowth()

	fmt.Println("\n=== 4. 硬核底层：Go 1.20+ 官方零拷贝转换 ===")
	demoZeroCopyConversion()
}

func demoBasicSliceAndMap() {
	// 切片 make 预分配容量
	users := make([]string, 0, 4)
	users = append(users, "Alice", "Bob")
	fmt.Printf("users 切片 len: %d, cap: %d, 内容: %v\n", len(users), cap(users), users)

	// Map comma-ok 惯用写法
	scores := map[string]int{"Alice": 95}
	if score, exists := scores["Bob"]; exists {
		fmt.Printf("Bob 分数: %d\n", score)
	} else {
		fmt.Println("Bob 成绩不存在 (返回 int 零值 0)")
	}
}

func demoFullSliceExpression() {
	// 1. 危险复现：普通切片污染原数组
	arrDangerous := []int{10, 20, 30, 40, 50}
	subDangerous := arrDangerous[1:3] // len=2, cap=4 (共享底层数组到末尾)
	subDangerous = append(subDangerous, 999) // 悄然将 arrDangerous[3] 覆盖为 999！
	fmt.Printf("危险切片 append 后，原切片被污染: %v\n", arrDangerous)

	// 2. 生产正解：三索引切片锁定 cap 隔离内存
	arrSafe := []int{10, 20, 30, 40, 50}
	subSafe := arrSafe[1:3:3] // 限制 cap = 3-1 = 2
	subSafe = append(subSafe, 999) // 容量已满，强行分配新数组
	fmt.Printf("三索引切片 append 后，原切片毫发无损: %v (子切片: %v)\n", arrSafe, subSafe)
}

func demoSliceGrowth() {
	var s []int
	var prevCap int

	fmt.Print("切片动态扩容 cap 跃迁轨迹: ")
	for i := 0; i < 2000; i++ {
		s = append(s, i)
		currentCap := cap(s)
		if currentCap != prevCap {
			fmt.Printf("%d ", currentCap)
			prevCap = currentCap
		}
	}
	fmt.Println("\n>> 证明: 1.18+ 平滑扩容曲线与内存对齐块分配，消除旧版断崖式跳跃。")
}

func demoZeroCopyConversion() {
	originalStr := "高性能后端研发"

	// 零拷贝 String -> []byte
	byteSlice := unsafe.Slice(unsafe.StringData(originalStr), len(originalStr))
	// 零拷贝 []byte -> String
	reconstructedStr := unsafe.String(unsafe.SliceData(byteSlice), len(byteSlice))

	fmt.Printf("原始字符串: %s\n", originalStr)
	fmt.Printf("零拷贝切片 len: %d, cap: %d\n", len(byteSlice), cap(byteSlice))
	fmt.Printf("零拷贝还原字符串: %s\n", reconstructedStr)
	fmt.Printf("原字符串数据地址: %p, 切片底层地址: %p (完全一致证明 0 拷贝！)\n",
		unsafe.StringData(originalStr), unsafe.SliceData(byteSlice))
}
