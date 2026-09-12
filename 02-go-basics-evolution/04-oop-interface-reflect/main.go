package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Counter struct {
	val int
}

// 值接收者：无法改变外部值
func (c Counter) AddCopy() {
	c.val++
}

// 指针接收者：能够真实改变外部值
func (c *Counter) AddReal() {
	c.val++
}

type CustomError struct {
	Code int
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("custom error code: %d", c.Code)
}

// 模拟容易出 bug 的错误返回函数
func returnTypedNilError() error {
	var err *CustomError = nil
	return err // 致命陷阱：返回了包装了类型指针的接口
}

func main() {
	fmt.Println("=== 1. 新手入门：值接收者 vs 指针接收者 ===")
	demoReceiverDifferences()

	fmt.Println("\n=== 2. 硬核底层：致命的 nil 接口判断陷阱 ===")
	demoNilInterfaceTrap()

	fmt.Println("\n=== 3. 硬核底层：iface 底层结构的内存拆解 ===")
	demoIfaceDissection()

	fmt.Println("\n=== 4. 硬核底层：反射基本机制与 CanSet 校验 ===")
	demoReflectionCanSet()
}

func demoReceiverDifferences() {
	c := Counter{val: 10}
	c.AddCopy()
	fmt.Printf("调用 AddCopy 后 val = %d (预期: 10，未被修改)\n", c.val)
	c.AddReal()
	fmt.Printf("调用 AddReal 后 val = %d (预期: 11，成功修改)\n", c.val)
}

func demoNilInterfaceTrap() {
	err := returnTypedNilError()

	// 绝大多数初学者会在此写出严重 Bug：
	if err == nil {
		fmt.Println("没有发生错误")
	} else {
		fmt.Println("【警告！发生错误！】即使内部指针是 nil，接口本身也不是 nil！")
		fmt.Printf("  -> 接口类型: %T, 接口值: %v\n", err, err)
	}
}

func demoIfaceDissection() {
	var err error = &CustomError{Code: 500}

	// 强制将 error 接口变量转换为模拟的 iface 结构体
	type dummyIface struct {
		tab  unsafe.Pointer
		data unsafe.Pointer
	}

	ifaceObj := *(*dummyIface)(unsafe.Pointer(&err))
	fmt.Printf("iface 内部指针剖析 -> tab(方法与类型表地址): %p, data(真实数据地址): %p\n",
		ifaceObj.tab, ifaceObj.data)
	fmt.Println(">> 证明: 只有当 tab 和 data 两个字段均为 0x0(nil) 时，interface == nil 才成立！")
}

func demoReflectionCanSet() {
	x := 42
	// 传递指针才能获取可寻址反射对象
	v := reflect.ValueOf(&x).Elem()

	if v.CanSet() {
		v.SetInt(100)
		fmt.Printf("通过反射成功修改值: x = %d\n", x)
	}
}
