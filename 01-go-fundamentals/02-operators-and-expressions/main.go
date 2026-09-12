package main

import (
	"fmt"
)

// 定义基于二进制位移的细粒度权限常量
const (
	PermRead    = 1 << 0 // 0001 (1)
	PermWrite   = 1 << 1 // 0010 (2)
	PermExecute = 1 << 2 // 0100 (4)
	PermAdmin   = 1 << 3 // 1000 (8)
)

// ClearPermission 演示使用 &^ 位清空运算符安全剥离权限
func ClearPermission(currentPerm int, permToRevoke int) int {
	return currentPerm &^ permToRevoke
}

// HasPermission 判定是否拥有指定权限
func HasPermission(currentPerm int, targetPerm int) bool {
	return (currentPerm & targetPerm) == targetPerm
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("     ⚡ 专题 02：运算符全景与位清空 (&^) 实操      ")
	fmt.Println("==================================================")

	// 1. 基础算术运算
	a, b := 17, 5
	fmt.Printf("算术运算: %d + %d = %d\n", a, b, a+b)
	fmt.Printf("整数截断除法: %d / %d = %d\n", a, b, a/b)
	fmt.Printf("求模取余数: %d %% %d = %d\n", a, b, a%b)

	// 2. 自增自减
	count := 10
	count++ // 合法独立语句
	fmt.Printf("自增后结果: %d (注意: Go 严禁使用 ++count 或赋值 x = count++)\n", count)

	// 3. RBAC 权限位清空实战
	userPerm := PermRead | PermWrite // 赋予读、写权限 (0011 = 3)
	fmt.Printf("\n--- RBAC 权限系统演练 ---\n")
	fmt.Printf("初始用户权限码: %04b, 包含写权限: %v\n", userPerm, HasPermission(userPerm, PermWrite))

	// 使用位清空运算符剥离写权限
	userPerm = ClearPermission(userPerm, PermWrite)
	fmt.Printf("剥离写权限后:   %04b, 包含写权限: %v, 包含读权限: %v\n",
		userPerm, HasPermission(userPerm, PermWrite), HasPermission(userPerm, PermRead))

	fmt.Println("==================================================")
}
