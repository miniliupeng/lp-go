package main

import (
	"fmt"
)

// User 演示基础结构体
type User struct {
	ID       int64
	Username string
	Age      int
	IsActive bool
}

// PrintGreeting 值接收者：仅操作局部拷贝
func (u User) PrintGreeting() {
	fmt.Printf("Hello, 我是 %s，今年 %d 岁\n", u.Username, u.Age)
}

// CelebrateBirthday 指针接收者：真实修改物理原对象
func (u *User) CelebrateBirthday() {
	u.Age++
}

// ------------------------------------------------------------
// 接口定义与多态实操
// ------------------------------------------------------------

// Engine 基础抽象接口
type Engine interface {
	RemainingMiles() int
	EngineType() string
}

// GasEngine 燃油引擎
type GasEngine struct {
	Gallons int
	MPG     int
}

func (g GasEngine) RemainingMiles() int {
	return g.Gallons * g.MPG
}

func (g GasEngine) EngineType() string {
	return "燃油引擎"
}

// ElectricEngine 电动引擎
type ElectricEngine struct {
	KWh   int
	MPKWh int
}

func (e ElectricEngine) RemainingMiles() int {
	return e.KWh * e.MPKWh
}

func (e ElectricEngine) EngineType() string {
	return "纯电引擎"
}

// CheckJourneyPlan 多态函数：接收任何具备 Engine 能力的对象
func CheckJourneyPlan(e Engine, distance int) bool {
	miles := e.RemainingMiles()
	canReach := miles >= distance
	fmt.Printf("【%s】总续航: %d 英里，前往 %d 英里目的地 -> 可达性: %t\n",
		e.EngineType(), miles, distance, canReach)
	return canReach
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("   🧩 专题 07：结构体、方法接收者与接口多态实操   ")
	fmt.Println("==================================================")

	// 1. 结构体初始化与接收者方法
	u := User{ID: 1001, Username: "张三", Age: 18, IsActive: true}
	fmt.Printf("1. 初始状态: %+v\n", u)

	u.PrintGreeting()
	u.CelebrateBirthday()
	fmt.Printf("   指针接收者方法调用后，Age 真实增加: %d\n", u.Age)

	// 2. 鸭子类型与接口多态
	fmt.Println("\n2. 接口非侵入式多态实测 (计划行程 250 英里):")
	gas := GasEngine{Gallons: 10, MPG: 30}     // 300 英里
	elec := ElectricEngine{KWh: 50, MPKWh: 4} // 200 英里

	CheckJourneyPlan(gas, 250)
	CheckJourneyPlan(elec, 250)

	fmt.Println("==================================================")
}
