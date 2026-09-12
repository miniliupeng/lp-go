package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// 1. 类型定义 (Type Definition): 创建全新类型，拥有独立方法集，不可直接混合赋值
type UserID int64

// 2. 类型别名 (Type Alias): 只是现有类型的别名，完全等价于原类型 (如 byte = uint8)
type ScoreAlias = int

// 3. 常量与高级 iota: 计算机存储容量单位计算
const (
	_  = iota             // 跳过 0
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB = 1 << (10 * iota) // 1 << 20
	GB = 1 << (10 * iota) // 1 << 30
)

// 4. 位掩码权限管理 (使用 Go 特有位运算符)
const (
	PermissionRead    = 1 << 0 // 0001 (1)
	PermissionWrite   = 1 << 1 // 0010 (2)
	PermissionExecute = 1 << 2 // 0100 (4)
)

// 结构体内存对齐验证结构体
type BadStruct struct {
	A bool  // 1 字节，偏移 0，后填 7 字节
	B int64 // 8 字节，偏移 8
	C bool  // 1 字节，偏移 16，后填 7 字节
}

type GoodStruct struct {
	B int64 // 8 字节，偏移 0
	A bool  // 1 字节，偏移 8
	C bool  // 1 字节，偏移 9，尾部填充 6 字节
}

func main() {
	fmt.Println("=== 1. 变量、零值与强类型转换 ===")
	demoVariablesAndConversions()

	fmt.Println("\n=== 2. make vs new 的天壤之别 ===")
	demoMakeVsNew()

	fmt.Println("\n=== 3. 类型定义 vs 类型别名 ===")
	demoTypeDefinitionVsAlias()

	fmt.Println("\n=== 4. iota 高级单位与 Go 特有位清空运算符 (&^) ===")
	demoAdvancedIotaAndBitOps()

	fmt.Println("\n=== 5. 字符串 byte vs rune 与高效拼接 ===")
	demoStringsAndRune()

	fmt.Println("\n=== 6. 切片 copy 深拷贝与 strconv 转换 ===")
	demoSliceCopyAndStrconv()

	fmt.Println("\n=== 7. 硬核底层：结构体内存对齐与 unsafe 穿透 ===")
	demoMemoryAlignmentAndUnsafe()
}

func demoVariablesAndConversions() {
	var num int
	var flag bool
	var text string
	var ptr *int

	fmt.Printf("零值安全体系: int=%d, bool=%t, string=%q, ptr=%v\n", num, flag, text, ptr)

	// Go 严禁隐式转换，哪怕 int 和 int64 也必须显式强转
	var a int = 100
	var b int64 = int64(a)
	fmt.Printf("显式类型转换: int(%d) -> int64(%d)\n", a, b)
}

func demoMakeVsNew() {
	// new(T): 分配内存，置零，返回指针 *T
	pInt := new(int)
	fmt.Printf("new(int) 返回指针: %p, 指向的值为零值: %d\n", pInt, *pInt)

	// ❌ 错误示范：如果用 new 创建 map，只能得到一个指向 nil map 的指针，写入直接 panic！
	// pMap := new(map[string]int)
	// (*pMap)["key"] = 1 // 💥 会 panic: assignment to entry in nil map

	// make: 专门用于 slice, map, channel 的内部数据结构初始化，返回实例本身
	m := make(map[string]int, 8)
	m["key"] = 100
	fmt.Printf("make 初始化 map 成功，len: %d, value: %d\n", len(m), m["key"])
}

func demoTypeDefinitionVsAlias() {
	var uid UserID = 10001
	var normalInt int64 = 10001

	// uid = normalInt // ❌ 编译报错！UserID 是全新类型，与 int64 无法直接赋值
	uid = UserID(normalInt) // 必须显式类型转换
	fmt.Printf("类型定义 UserID: %v, 类型名: %T\n", uid, uid)

	var score ScoreAlias = 99
	var realInt int = score // ✅ 完全等价，可以直接赋值给 int
	fmt.Printf("类型别名 ScoreAlias: %v, 真实底层类型: %T\n", realInt, score)
}

func demoAdvancedIotaAndBitOps() {
	fmt.Printf("存储容量常量: 1KB=%d 字节, 1MB=%d 字节, 1GB=%d 字节\n", KB, MB, GB)

	// 权限位掩码实操
	userPerm := PermissionRead | PermissionWrite // 拥有读写权限 (0011)
	hasWrite := (userPerm & PermissionWrite) != 0
	fmt.Printf("用户初始权限: 拥有写权限? %t\n", hasWrite)

	// Go 特有位清空运算符 (&^: AND NOT): 将指定位置清零
	userPerm = userPerm &^ PermissionWrite // 清除写权限
	hasWriteAfter := (userPerm & PermissionWrite) != 0
	fmt.Printf("使用 &^ 清除写权限后: 拥有写权限? %t\n", hasWriteAfter)
}

func demoStringsAndRune() {
	msg := "Go语言"
	fmt.Printf("字符串 %q 底层字节数 len(): %d, 真实字符数: %d\n",
		msg, len(msg), utf8.RuneCountInString(msg))

	// strings.Builder 高效拼接
	var builder strings.Builder
	builder.Grow(32) // 预分配内存，彻底消除扩容拷贝
	builder.WriteString("Hello")
	builder.WriteString(" ")
	builder.WriteString("Gopher")
	fmt.Printf("strings.Builder 拼接结果: %s\n", builder.String())
}

func demoSliceCopyAndStrconv() {
	src := []int{1, 2, 3, 4, 5}
	// 新手巨坑：copy 前目标切片必须具备足够 len，若是 make([]int, 0) 则复制 0 个元素！
	dst := make([]int, len(src))
	copy(dst, src)
	dst[0] = 999
	fmt.Printf("copy 深拷贝后原切片未受影响: src[0]=%d, dst[0]=%d\n", src[0], dst[0])

	// strconv 类型转换
	valStr := "2026"
	valInt, _ := strconv.Atoi(valStr)
	fmt.Printf("strconv.Atoi: %q -> %d (自增结果: %d)\n", valStr, valInt, valInt+1)
}

func demoMemoryAlignmentAndUnsafe() {
	var bad BadStruct
	var good GoodStruct

	fmt.Printf("BadStruct  尺寸: %d 字节, GoodStruct (重排优化) 尺寸: %d 字节\n",
		unsafe.Sizeof(bad), unsafe.Sizeof(good))
	fmt.Printf(">> 结构体字段重排立省 %.1f%% 内存！\n",
		float64(unsafe.Sizeof(bad)-unsafe.Sizeof(good))/float64(unsafe.Sizeof(bad))*100)

	type Empty struct{}
	var e1, e2 Empty
	fmt.Printf("空结构体 struct{} 大小: %d 字节, e1 地址: %p, e2 地址: %p (统一指向 runtime.zerobase)\n",
		unsafe.Sizeof(e1), &e1, &e2)
}
