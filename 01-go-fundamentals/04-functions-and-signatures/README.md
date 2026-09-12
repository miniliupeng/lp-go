# 专题 04：函数签名、多返回值与闭包入门

> 函数是 Go 语言的**一等公民（First-Class Citizen）**：函数不仅能独立声明，还可以作为参数传入、作为返回值输出，或赋值给局部变量。Go 语言最标志性的语言特色之一就是原生支持**多返回值**，彻底革新了传统语言单一返回值的局限。

---

## 1. 函数声明语法与形参合并

```go
// 当连续多个参数类型相同时，可省略前面的类型，合并在最后一个参数后声明
func Add(a, b int) int {
    return a + b
}
```

---

## 2. 标志性特性：原生多返回值与命名返回值

在其他语言（如 Java/C++）中，如果一个函数既要返回计算结果、又要返回是否成功，往往需要封装一个对象或者借助输出参数；而在 Go 语言中，原生支持直接返回多个值：

```go
// 1. 普通多返回值（最常见于业务返回：结果 + error）
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 2. 命名返回值（在返回值列表中直接声明变量名）
// 此时返回值在函数入口处自动初始化为零值，结尾可直接使用裸 return (naked return)
func CalculateRect(width, height float64) (area float64, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return // 自动将 area 和 perimeter 拷贝输出
}
```

> [!TIP]
> 命名返回值能够极好地充当自解释文档；但在长函数中建议显式 `return area, perimeter`，避免裸 return 降低长代码的可读性。

---

## 3. 可变参数列表（Variadic Parameters）

使用 `...Type` 语法，函数可以接收任意数量的同类型参数，在函数内部该参数被当作切片（Slice）处理：

```go
func Sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 调用时可以传任意个数字，也可以直接解包现有切片：
Sum(1, 2, 3)
slice := []int{4, 5, 6}
Sum(slice...) // 使用 ... 解包切片传入
```

---

## 4. 匿名函数与闭包（Closures）入门

闭包由**一个函数**与其引用的**外部自由变量环境**共同组成。闭包函数“捕获”了外部变量的生命周期：

```go
// CreateCounter 返回一个闭包函数，它持久持有独立的 count 变量状态
func CreateCounter() func() int {
    count := 0
    return func() int {
        count++ // 闭包捕获了外部局部变量 count，即使 CreateCounter 已经返回，count 也依然存活！
        return count
    }
}
```

---

## 5. 参数传递哲学铁律：全盘值拷贝传递（Pass-by-Value）

Go 语言在函数传参时**永远是值拷贝（Pass-by-value）**！
- 如果你把一个普通变量 `x` 传给函数并在函数内修改，函数外部的 `x` **绝对不会改变**；
- 若想在函数内部修改外部变量的值，**必须显式传递该变量的指针（`*Type`）**！

---

## 6. 本地实操运行

```bash
go run ./01-go-fundamentals/04-functions-and-signatures/main.go
```
