package main

import (
	"fmt"
	"unsafe"
)

// 定义基于 iota 的业务状态枚举
const (
	OrderStateCreated = iota // 0
	OrderStatePaid           // 1
	OrderStateShipped        // 2
	OrderStateFinished       // 3
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("    📊 专题 01：Go 语言数据类型与变量声明演示     ")
	fmt.Println("==================================================")

	// 1. 基本数据类型声明与内存尺寸
	var aBool bool = true
	var anInt int = 42
	var anInt64 int64 = 9223372036854775807
	var aFloat float64 = 3.1415926535
	var aByte byte = 'A'
	var aRune rune = '中'
	var aString string = "Go语言2026"

	fmt.Println("--- 1. 数据类型实测与内存占用 (Bytes) ---")
	fmt.Printf("bool    : 值=%v, 尺寸=%d Byte\n", aBool, unsafe.Sizeof(aBool))
	fmt.Printf("int     : 值=%v, 尺寸=%d Bytes\n", anInt, unsafe.Sizeof(anInt))
	fmt.Printf("int64   : 值=%v, 尺寸=%d Bytes\n", anInt64, unsafe.Sizeof(anInt64))
	fmt.Printf("float64 : 值=%v, 尺寸=%d Bytes\n", aFloat, unsafe.Sizeof(aFloat))
	fmt.Printf("byte    : 字符='%c', ASCII码=%d, 尺寸=%d Byte\n", aByte, aByte, unsafe.Sizeof(aByte))
	fmt.Printf("rune    : 字符='%c', Unicode码=%d, 尺寸=%d Bytes\n", aRune, aRune, unsafe.Sizeof(aRune))
	fmt.Printf("string  : 值=\"%s\", 内部结构尺寸=%d Bytes\n", aString, unsafe.Sizeof(aString))

	// 反引号原生字符串字面量（Raw String Literal）：保留换行与特殊字符，无需转义
	rawString := `Line 1: SELECT * FROM users
Line 2: WHERE id = "1001";`
	fmt.Printf("rawString:\n%s\n", rawString)

	// 2. 零值体系（未赋值自动初始化为安全的默认零值）
	var defaultInt int
	var defaultBool bool
	var defaultString string
	var defaultPtr *int

	fmt.Println("\n--- 2. 默认零值体系（绝无脏内存数据）---")
	fmt.Printf("未赋值 int 零值    : %d\n", defaultInt)
	fmt.Printf("未赋值 bool 零值   : %v\n", defaultBool)
	fmt.Printf("未赋值 string 零值 : \"%s\" (长度=%d)\n", defaultString, len(defaultString))
	fmt.Printf("未赋值 指针 零值   : %v\n", defaultPtr)

	// 3. iota 枚举输出
	fmt.Println("\n--- 3. iota 优雅自增枚举 ---")
	fmt.Printf("订单状态: Created=%d, Paid=%d, Shipped=%d, Finished=%d\n",
		OrderStateCreated, OrderStatePaid, OrderStateShipped, OrderStateFinished)
	fmt.Println("==================================================")
}
