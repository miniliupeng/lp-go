# 04 控制流

## 这一章你要会什么

用 `if`、`for`、`switch` 写清楚分支和循环。Go 没有 `while`，也没有三元运算符。

## if

条件和语句体之间没有括号。可以带一条只在 `if` / `else` 里有效的短声明，很常见的写法是 `if v, err := fn(); err != nil`。

## for

三种都是 `for`：

```go
for i := 1; i <= 3; i++ { /* 计数 */ }
for condition { /* 相当于 while */ }
for { /* 无限循环，用 return 或 break 离开 */ }
```

`for i, v := range slice` 遍历切片、数组、字符串、map、channel。

**Go 1.22 起，每次迭代的循环变量都是新的。** `for _, n := range nums { go func() { use(n) }() }` 在 1.22 及以后会捕获当次的 `n`。更早的版本会把所有 goroutine 都绑到同一个变量上。本教程的并发示例仍然把 `n` 当参数传进去，这样读代码的人不用先回忆语言版本。

Go 1.22 起还可以 `for i := range n`，`i` 从 0 到 `n-1`。`n` 是整数。

`break` 离开当前循环，`continue` 进入下一次。想离开外层循环，给循环加标签。

## switch

`switch` 命中一个 `case` 后就结束，**不会**掉进下一个 `case`。多个值写成 `case 2, 3:`。需要故意继续时才写 `fallthrough`。`switch` 也可以没有操作数，每个 `case` 写一个条件。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/control
```

输出：

```text
odd
double=6
sum=6
rangeN=[0 1 2]
two-or-three
```

## 练习

把 `switch` 改成判断一个字符串命令：`"sum"` 时做加法，其他值走 `default`。
