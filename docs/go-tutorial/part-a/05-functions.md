# 05 函数

## 这一章你要会什么

写多返回值函数、可变参数、函数值和 `defer`。

## 多返回值

Go 用额外的返回值表示失败，而不是抛异常。习惯是最后一个返回值的类型是 `error`。

```go
func div(a, b int) (int, int, error) {
    if b == 0 {
        return 0, 0, errors.New("divide by zero")
    }
    return a / b, a % b, nil
}
```

可以给返回值起名字。名字会变成函数里的局部变量，裸 `return` 会把它们返回去。裸 `return` 在长函数里难读，短函数里偶尔有用。本教程的示例都把返回值写出来。

调用时如果不要某个结果，用 `_`。

## 可变参数

`func sum(nums ...int)` 在函数内部就是 `[]int`。调用时 `sum(1, 2, 3)`。已经有切片时用 `sum(s...)`。

## 函数值和闭包

函数可以赋给变量，也可以当参数。闭包会捕获外部变量。捕获的是变量本身，不是捕获那一瞬间的拷贝。所以计数器可以这样写：

```go
func counter() func() int {
    n := 0
    return func() int {
        n++
        return n
    }
}
```

## defer

`defer` 把一次函数调用推迟到当前函数返回前。多个 `defer` 按**后进先出**执行。适合配对的资源：`f, err := os.Open(path)` 成功后立刻 `defer f.Close()`。

`defer` 的参数在执行到 `defer` 这一行时就求值，不是等到函数返回时再求值。

## 参数是复制

`func add(a, b int)` 得到的是副本。要修改调用方的结构体，用指针接收者或指针参数。复制一个切片头不会复制底层数组。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/functions
```

输出：

```text
div=3 rem=1 err=<nil>
divZero=divide by zero
sum=10
closure=1
closure=2
body
defer 2
defer 1
```

`err=<nil>` 是 `fmt` 打印 nil 错误的方式，表示没有错误。

## 练习

写 `minMax(nums ...int) (int, int, error)`。`nums` 为空时返回错误。用 `defer` 打印 `"done"`，确认它出现在结果之后。
