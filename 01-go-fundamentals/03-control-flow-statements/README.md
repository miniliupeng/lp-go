# 专题 03：流程控制全景：条件分支、四态循环与 Switch

> 任何程序的逻辑骨架都离不开流程控制。Go 语言在控制流设计上遵循“极简与安全”哲学：**全语言只有一种循环关键字（`for`）**，彻底取缔了 `while` 与 `do-while`；同时 `switch` 默认自带 `break`，防止遗忘导致的向下贯通 bug。

---

## 1. 条件分支：`if - else` 与特有的前置短声明

Go 的 `if` 条件无需加小括号 `()`，且允许在条件判断前执行一段**短变量初始化语句**：

```go
// 经典规范：初始化与判定合二为一
// 变量 err 的作用域严格限制在 if-else 块内部，外部无法访问，极大防止局部变量污染！
if err := executeTask(); err != nil {
    fmt.Printf("任务执行失败: %v\n", err)
    return
}
```

### 1.1 ⚠️ 语法禁忌：`else if` / `else` 换行报错与自动分号机制（ASI）

许多从 C++、Java 或 Python 转来的开发者，容易习惯性将 `else` 单独换行写在新的一行，但这在 Go 语言中是**绝对被禁止的编译错误**：

::: code-group
```go [❌ 错误写法：换行导致编译失败]
if score >= 90 {
    fmt.Println("优秀")
}
// 编译器直接报错：syntax error: unexpected else, expected statement
else {
    fmt.Println("及格")
}
```

```go [✅ 正确写法：紧随闭合大括号同行]
if score >= 90 {
    fmt.Println("优秀")
} else if score >= 60 {
    fmt.Println("及格")
} else {
    fmt.Println("不及格")
}
```
:::

#### 🔬 底层第一性原理：词法分析器的“自动分号插入”（ASI）
为什么 Go 编译器如此严苛？其根源在于 **Go 语言规范中的自动分号插入机制（Automatic Semicolon Insertion）**：
1. 为了保持代码清爽，Go 允许开发者在每行代码末尾省略分号 `;`；
2. Go 词法分析器（Lexer）在词法扫描阶段，一旦遇到单独成行的闭合大括号 `}`、关键字 `return`、`break` 或标识符，**会自动在该行末尾隐式补上一个分号 `;`**；
3. 如果将 `else` 写在下一行，前一行的 `}` 会被自动补上分号变为：
   ```go
   if score >= 90 {
       fmt.Println("优秀")
   }; // <-- 词法分析器自动补全分号，判定 if 语句在此合法终结！
   else { // <-- 变成一个突兀且孤立的 else，词法解析器报错！
   ```
4. 因此，将 `} else {` 连写同行是 Go 语法的硬性设计约束。

---

## 2. 循环体系：Go 只有唯一的 `for` 关键字（四大形态）

Go 没有 `while` 和 `do-while`，所有的循环均通过 `for` 实现：

### ① 经典三段式循环
```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

### ② 单条件循环（等价于其他语言的 `while(condition)`）
```go
n := 1
for n < 100 {
    n *= 2
}
```

### ③ 无限死循环（等价于 `while(true)`）
```go
for {
    // 配合 select、Channel 或特定退出条件 break
    if shouldExit {
        break
    }
}
```

### ④ 现代 Go 1.22+ 整数遍历语法糖
```go
// 直接遍历 0 到 9，无需繁琐的三段式声明！
for i := range 10 {
    fmt.Println(i) // 打印 0, 1, ..., 9
}
```

---

## 3. 多路选择：`switch - case` 的三大现代特性

1. **默认自动中断（无需写 `break`）**：
   - 匹配到一个 `case` 并执行完毕后，程序自动跳出 `switch`。
   - 若刻意需要贯通执行下一个 case，需显式使用 **`fallthrough`** 关键字。
2. **支持多值匹配**：
   - `case "A", "B", "C":`（只要满足其中一个即可触发）。
3. **无表达式的 `switch`（替代繁琐的 `if-else if-else`）**：
   ```go
   score := 85
   switch {
   case score >= 90:
       fmt.Println("优秀")
   case score >= 80:
       fmt.Println("良好")
   default:
       fmt.Println("合格")
   }
   ```

---

## 4. 标签（Label）与跳出外层多重循环

在多层嵌套循环中，普通 `break` 只能跳出当前最内层循环；配合 **Label 标签**，可以一键精准跳出外层指定循环：

```go
OuterLoop:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if i*j > 6 {
            break OuterLoop // 直接终结外层 OuterLoop 循环！
        }
    }
}
```

---

## 5. 本地实操运行

```bash
go run ./01-go-fundamentals/03-control-flow-statements/main.go
```
