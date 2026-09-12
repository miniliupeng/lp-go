# 专题 07：结构体定义、方法接收者与非侵入式接口多态实战

> 在传统面向对象语言（如 Java、C++）中，类（`class`）与继承（`extends`）是构建系统抽象的核心；而在 Go 语言中，完全**取缔了 `class` 和继承**，转而采用更加轻盈正交的 **“结构体（struct）+ 方法接收者（Receiver）+ 鸭子类型接口（Interface）”**。本专题带你从零掌握面向对象在 Go 中的正确打开方式。

---

## 1. 结构体（`struct`）声明与初始化姿势

结构体是用户自定义的复合类型，将多个不同类型的字段组合在一起：

```go
// 1. 声明结构体
type User struct {
    ID       int64  // 唯一 ID
    Username string // 用户名
    Age      int    // 年龄
    IsActive bool   // 是否激活
}

// 2. 初始化姿势 A：字段名键值对声明（大厂生产推荐！即使结构体增减字段也不会报错）
u1 := User{
    ID:       1001,
    Username: "Gopher",
    Age:      18,
    IsActive: true,
}

// 3. 初始化姿势 B：按字段顺序省略键名（不推荐，容易错位）
u2 := User{1002, "Alice", 20, false}

// 4. 初始化姿势 C：使用 new 返回指针（分配内存并赋予默认零值）
uPtr := new(User) // 等价于 &User{}，uPtr 为 *User 类型
```

---

## 2. 方法接收者：值接收者 vs 指针接收者（深层剖析）

在 Go 中，方法（Method）是**附带了“接收者形参（Receiver）”的特殊函数**：

```go
// 姿势 1：值接收者 (Value Receiver)
// 传入的是结构体的一份局部副本（内存拷贝），内部修改对外部原对象完全无影响！
func (u User) PrintGreeting() {
    fmt.Printf("Hello, 我是 %s，今年 %d 岁\n", u.Username, u.Age)
}

// 姿势 2：指针接收者 (Pointer Receiver)
// 传入的是物理内存地址，内部修改会真实映射到外部原对象上！
func (u *User) UpdateAge(newAge int) {
    u.Age = newAge // 自动解引用并真实修改
}
```

### 🔬 什么时候选值接收者？什么时候选指针接收者？

| 考量维度 | 值接收者 (`u User`) | 指针接收者 (`u *User`) |
| :--- | :--- | :--- |
| **是否修改原对象** | **仅读取**，不修改字段 | **必须修改**内部字段或状态 |
| **内存开销** | 每次调用产生结构体**全量字节拷贝** | 仅拷贝 8 字节物理内存地址指针 |
| **包含锁字段** | ❌ 严禁使用（拷贝会导致锁失效） | ✅ **必须使用**（保证所有方法操作同一把锁） |
| **大厂工程规范** | 小型轻量对象（如包含 1~2 个字段的坐标点） | **大厂 95% 场景默认首选指针接收者** |

---

## 3. 非侵入式接口（Interface）与鸭子类型（Duck Typing）

在 Java 中实现接口必须写 `class Dog implements Animal`；而在 Go 语言中采用**隐式实现（Duck Typing）**：
> *“如果它走起路来像鸭子，叫起来也像鸭子，那么它就是一只鸭子。”*

只要一个结构体拥有了接口声明的**全部方法**，Go 编译器就会在编译期自动判定该类型实现了此接口，**零关键字显式声明**！

### 3.1 定义抽象契约与具体实现

```go
// 1. 定义抽象接口：定义一组“行为能力契约”
type Engine interface {
    RemainingMiles() int // 契约：必须具备计算剩余里程的能力
}

// 2. 具体类型 A：燃油引擎结构体
type GasEngine struct {
    Gallons int // 剩余油量（加仑）
    MPG     int // 每加仑行驶里程
}

// 隐式实现 Engine 接口（拥有 RemainingMiles 方法）
func (g GasEngine) RemainingMiles() int {
    return g.Gallons * g.MPG
}

// 3. 具体类型 B：电动引擎结构体
type ElectricEngine struct {
    KWh   int // 剩余电池电量（度）
    MPKWh int // 每度电行驶里程
}

// 隐式实现 Engine 接口
func (e ElectricEngine) RemainingMiles() int {
    return e.KWh * e.MPKWh
}
```

### 3.2 接口多态：统一调用解耦底层拓扑

通过将接口作为函数形参，调用方无需关心传入的是燃油车还是电动车，彻底实现代码解耦与面向契约编程：

```go
// 多态函数：接收任何具备 Engine 能力的对象
func CanReachDestination(e Engine, distance int) bool {
    return e.RemainingMiles() >= distance
}

func main() {
    gasCar := GasEngine{Gallons: 10, MPG: 30}          // 续航 300 英里
    elecCar := ElectricEngine{KWh: 60, MPKWh: 4}      // 续航 240 英里

    fmt.Println(CanReachDestination(gasCar, 250))  // true
    fmt.Println(CanReachDestination(elecCar, 250)) // false
}
```

> [!NOTE]
> **平滑进阶指引**：掌握上述非侵入式契约与多态传参已能覆盖 90% 的业务场景。关于接口在内存底层的动态类型信息、虚函数调用表以及 `iface/eface` 结构体原理，我们在 **阶段二「02-04 接口进阶与底层脱糖」** 中展开深度剖析。

---

## 4. 本地实操运行与测试验证

```bash
# 运行结构体与接口多态演示
go run ./01-go-fundamentals/07-structs-and-interfaces/main.go

# 运行自动化测试
go test -v ./01-go-fundamentals/07-structs-and-interfaces/...
```
