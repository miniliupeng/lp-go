# 专题 06：指针入门与数组/切片/字典基础 CRUD 实战

> 本专题带领新手平缓跨越 Go 语言开发中最常用的三个核心构件：**指针（Pointer）**、**切片（Slice）** 和 **字典（Map）**。本章聚焦日常业务开发必须掌握的增删改查（CRUD）基础用法，建立清晰实用的开发心智。

---

## 1. 指针（Pointer）极速入门：`&` 与 `*`

在 Go 中，指针非常安全（默认不支持像 C 语言那样的指针算术运算）：
- **`&`（取地址符）**：获取某个变量在物理内存中的十六进制地址（例如 `&a`）；
- **`*`（解引用符）**：通过内存地址，访问或修改该地址上存储的真实变量值（例如 `*ptr = 100`）。

```go
func DoubleInPlace(ptr *int) {
    *ptr = *ptr * 2 // 直接在原内存位置将数值翻倍
}

x := 10
DoubleInPlace(&x) // 传入 x 的内存地址
// 此时外部的 x 变成了 20！
```

---

## 2. 数组（Array）vs 切片（Slice）

| 特性维度 | 数组（Array） | 切片（Slice，最常用！） |
| :--- | :--- | :--- |
| **类型声明** | `[3]int`（方括号内有**固定长度**） | `[]int`（方括号内**无数字**，长度可变） |
| **推导语法** | `[...]int{1, 2, 3}`（由编译器自动推导长度） | `make([]int, 0, 5)` 或 `[]int{1, 2, 3}` |
| **传参行为** | **纯值类型**，传参会整体验证并进行全量内存拷贝 | **引用视图**，内部持有底层数组指针，性能开销极低 |
| **动态扩容** | 不可扩容，尺寸编译期定死 | 原生支持 `append` 动态追加自动扩容 |

```go
// 显式指定长度
var arr1 [3]int = [3]int{1, 2, 3}

// 编译器自动推导长度（仍是固定长度数组，类型为 [3]int）
arr2 := [...]int{1, 2, 3}
```

---

## 3. 切片（Slice）的四大日常操作

```go
// 1. 创建切片（推荐使用 make 指定初始长度与容量）
nums := make([]int, 0, 5) // len=0, cap=5

// 2. 动态追加元素（append）
nums = append(nums, 10, 20, 30)

// 3. 切片截取操作 [start:end]（包含 start，不包含 end）
sub := nums[1:3] // 截取索引 1 和 2，得到 [20, 30]

// 4. 遍历切片（for-range）
for idx, val := range nums {
    fmt.Printf("索引=%d, 数值=%d\n", idx, val)
}
```

---

## 4. 字典（Map）日常基础增删改查（CRUD）

Map 是键值对（Key-Value）哈希表。**必须使用 `make` 初始化后才能写入数据！**

```go
// 1. 初始化（如果未初始化为 nil，直接写 map 会触发 panic！）
userAges := make(map[string]int)

// 2. 写入与更新（Create / Update）
userAges["张三"] = 25
userAges["李四"] = 30

// 3. 读取与 comma-ok 安全判定（Read）
// 如果 key 存在，ok 为 true；如果 key 不存在，ok 为 false（且 val 返回默认零值）
if age, ok := userAges["张三"]; ok {
    fmt.Printf("张三的年龄是: %d\n", age)
} else {
    fmt.Println("未找到该用户")
}

// 4. 删除键值对（Delete）
delete(userAges, "李四") // 安全删除，若 key 不存在也不会报错
```

### 4.1 ⚠️ 为什么 `for-range` 遍历 Map 的顺序是完全随机的？

如果你写循环遍历同一个 Map：
```go
for k, v := range userAges {
    fmt.Println(k, v)
}
```
你会惊奇地发现：**每次程序运行，打印出来的键值对顺序可能都不一样！**

> [!IMPORTANT]
> **底层第一性原理：防脆弱随机哈希种子**
> * 在很多早期编程语言中，开发者会写出“隐式依赖 Map 默认插入顺序”的代码，一旦升级语言或重构哈希算法，代码立刻产生幽灵 bug；
> * Go 设计团队为了**从根源上彻底断绝这种坏习惯**，在运行时内部（`runtime.mapiterinit`）故意引入了**随机种子偏移（fastrand）**；
> * 哪怕哈希桶里的物理数据完全一样，每次遍历都会从随机的一个桶编号开始。这是典型的工业级**防脆弱（Anti-fragile）工程设计**！

---

## 5. 本地实操运行与验证

```bash
# 运行指针与切片/字典基础演示
go run ./01-go-fundamentals/06-pointers-and-slices-maps/main.go

# 运行本专题自动化单元测试
go test -v ./01-go-fundamentals/06-pointers-and-slices-maps/...
```

