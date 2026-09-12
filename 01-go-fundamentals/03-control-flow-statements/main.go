package main

import (
	"errors"
	"fmt"
)

// EvaluateScore 演示无条件表达式的 switch-case 判定
func EvaluateScore(score int) string {
	switch {
	case score >= 90:
		return "优秀 (Grade A)"
	case score >= 80:
		return "良好 (Grade B)"
	case score >= 60:
		return "及格 (Grade C)"
	default:
		return "未达标 (Grade D)"
	}
}

// FindFirstMatrixIntersection 演示使用 Label 标签优雅跳出多层嵌套循环
func FindFirstMatrixIntersection(matrix [][]int, target int) (int, int, bool) {
	found := false
	targetRow, targetCol := -1, -1

SearchLoop:
	for r, row := range matrix {
		for c, val := range row {
			if val == target {
				targetRow, targetCol = r, c
				found = true
				break SearchLoop // 直接跳出外层整个循环，无需多重 flag 标记！
			}
		}
	}

	return targetRow, targetCol, found
}

func mockCheckUser(id int) error {
	if id <= 0 {
		return errors.New("用户 ID 必须为正整数")
	}
	return nil
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("     🔀 专题 03：流程控制全景与多态循环实操       ")
	fmt.Println("==================================================")

	// 1. if 前置声明语句实战
	userID := -1
	if err := mockCheckUser(userID); err != nil {
		fmt.Printf("1. if 前置短声明成功拦截非法输入: %v\n", err)
	}

	// 2. 现代 Go 1.22+ 遍历演示
	fmt.Print("2. Go 1.22+ 现代 `for i := range 5` 序列输出: ")
	for i := range 5 {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// 3. switch 判定
	fmt.Printf("3. 分数等级评估 (85分) -> %s\n", EvaluateScore(85))

	// 4. 多层循环 Label 跳出演示
	matrix := [][]int{
		{1, 2, 3},
		{4, 99, 6},
		{7, 8, 9},
	}
	r, c, found := FindFirstMatrixIntersection(matrix, 99)
	fmt.Printf("4. 矩阵目标 99 命中坐标: 行=%d, 列=%d, 是否命中=%v (Label 直接跳出)\n", r, c, found)

	fmt.Println("==================================================")
}
