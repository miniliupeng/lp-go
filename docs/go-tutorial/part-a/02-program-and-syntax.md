# 02 程序结构与声明

## 这一章你要会什么

看懂一个 Go 文件由什么组成，会声明变量和常量，知道零值是什么。

## 文件里有什么

一个文件以 `package` 开头。同一个目录里的文件必须是同一个包名。`import` 引入其他包。`main` 包里的 `main` 函数是入口。

Go 用大写开头表示**对其他包可见**。`Add` 可以被别的包调用，`add` 不行。这不是关键字，是语言规则。

语句末尾的分号由编译器插入。`if` 的 `{` 必须和 `if` 写在同一行，不能换到下一行。把代码交给 `go fmt`，不要手调对齐。

注释用 `//` 或 `/* */`。

## 变量

```go
var language string
language = "Go"
year := 2009
```

`var` 可以只声明。`:=` 是短声明，只能写在函数内部，并且左值里至少要有一个新变量。

函数外不能用 `:=`。包一级的变量用 `var`。

## 常量与 iota

```go
const (
    Red int = iota // 0
    Green          // 1
    Blue           // 2
)
```

`iota` 在每个 `const` 块从 0 开始，每下一行加 1。适合一组相关的整数常量。它不是枚举类型系统，没有自动的类型安全。

## 零值

声明了但没赋值的变量是零值，不是垃圾内存：

| 类型 | 零值 |
| --- | --- |
| 数值 | `0` |
| `bool` | `false` |
| `string` | `""` |
| 指针、切片、map、channel、接口、函数 | `nil` |

所以「声明一个字符串再读它」是合法的，得到的是空串，不会像未初始化的 C 局部变量那样是随机值。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/syntax
```

输出：

```text
language=Go year=2009 title=tutorial
zeroInt=0 zeroString="" zeroBool=false
Red=0 Green=1 Blue=2
sum=3
```

源码在 [examples/parta/syntax/main.go](../examples/parta/syntax/main.go)。

## 练习

1. 在 `const` 块里加一个 `Yellow`，预测它的值再运行。
2. 写一个 `var ready bool`，不要赋值，用 `fmt.Printf("%t", ready)` 看零值。
