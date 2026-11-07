# 04 控制流

## 这一章你要会什么

用 `if`、`for`、`switch` 写分支和循环，并说清 `break`、`continue`、`fallthrough` 各自离开的是哪一层。Go 没有 `while`，也没有三元运算符。

## if

条件和语句体之间没有括号。`if` 可以带一条短声明，这条声明只在这个 `if` 和配套的 `else` 里有效：

```go
if v := n * 2; v > 5 {
    fmt.Printf("double=%d\n", v)
}
```

`v` 出了这对括号就不能用。最常见的形式是 `if v, err := fn(); err != nil`。

## for

三种循环都写成 `for`：

```go
for i := 1; i <= 3; i++ { /* 计数。i++ 是语句，写在这里 */ }
for condition { /* 条件为真就继续，相当于别的语言里的 while */ }
for { /* 没有条件。用 return 或 break 离开 */ }
```

`for i, v := range x` 能遍历数组、切片、字符串、map 和 channel。遍历字符串时，`i` 是字节下标，`v` 是一个 `rune`，不是一个字节。遍历 channel 时只有一个变量，channel 被关闭并且剩余值取完后循环结束。

**Go 1.22 起，每次迭代的循环变量都是新的。** 在 1.22 及以后，`for _, n := range nums { go func() { use(n) }() }` 捕获的是当次的 `n`。更早的版本里这些 goroutine 会共享同一个 `n`。本教程的并发示例仍然把 `n` 当参数传进函数，读代码的人不用先回忆语言版本。

Go 1.22 起还可以 `for i := range n`，`n` 是整数，`i` 从 0 取到 `n-1`。示例里 `range 3` 得到 `0, 1, 2`。

`break` 离开当前这一层 `for`、`switch` 或 `select`。`continue` 进入当前这一层 `for` 的下一次迭代。要离开外层循环，给外层加上标签，再 `break` 那个标签。示例里内层走到 `a == 1 && b == 1` 时 `break outer`，外层不会继续跑到 `a == 2`，打印 `brokeOuter=true`。

## switch

`switch` 命中一个 `case` 之后就结束，不会掉进下一个 `case`。多个值写成 `case 2, 3:`。示例里 `n` 是 3，所以打印 `two-or-three`。

需要故意进入下一个 `case` 的语句时，当前 `case` 的最后一条写成 `fallthrough`。它只是跳到下一段语句，**不会**再计算下一个 `case` 的表达式。示例里 `n` 是 3，命中 `case 2, 3` 后 `fallthrough` 到 `default` 的打印，输出 `fall=through`。类型 `switch` 里不能写 `fallthrough`。

`switch` 也可以不写操作数。这时每个 `case` 是一个条件，从上到下试，第一个为真的分支执行。示例用它判断命令字符串，打印 `cmd=sum`。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/control
```

在本仓库的 Go 1.27.1 上，输出是：

```text
odd
double=6
sum=6
rangeN=[0 1 2]
two-or-three
cmd=sum
fall=through
brokeOuter=true
```

源码在 [examples/parta/control/main.go](../examples/parta/control/main.go)。

## 练习

1. 把 `cmd` 改成 `"other"`，预测两行新输出，再运行。
2. 去掉 `break outer` 的标签，只留 `break`，看 `brokeOuter` 会不会变成 `false`，以及外层还会不会继续。
