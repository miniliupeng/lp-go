# 03 类型

## 这一章你要会什么

会用基本类型、切片、map、字符串和指针，并且知道哪些转换必须写出来。

## 基本类型

常用的有 `bool`、`int`、`int64`、`float64`、`string`、`byte`、`rune`。

`byte` 是 `uint8` 的别名。`rune` 是 `int32` 的别名，用来表示一个 Unicode 码点。`int` 的宽度跟平台有关，不要用它保存必须是 32 位或 64 位的协议字段，那种地方写 `int32` 或 `int64`。

不同类型不能混在一个表达式里。`int` 和 `int64` 相加之前要显式转换：

```go
var n int = 65
s := string(rune(n)) // "A"，因为 U+0041 是 A
```

`string(65)` 在新代码里会被 vet 提醒：把整数直接转成字符串得到的是 Unicode 码点，不是十进制文本 `"65"`。十进制文本用 `strconv.Itoa`。

## 字符串、字节、rune

字符串是不可变的字节序列，通常按 UTF-8 解释。

`len("中")` 是 **3**（字节数），`len([]rune("中"))` 是 **1**。按下标 `s[i]` 得到的是字节，不是字符。逐个字符用 `for _, r := range s`。

## 数组和切片

数组长度是类型的一部分：`[3]int` 和 `[4]int` 是不同类型。数组赋值会复制全部元素。

切片是一段描述：指向底层数组、长度、容量。它是日常使用的动态序列。

```go
nums := []int{1, 2, 3}
nums = append(nums, 4)
```

`append` 在容量不够时分配新的底层数组。**容量够的时候，`append` 会改掉和你共享底层数组的其他切片。** 把切片传给函数时要记住这一点。

`make([]int, len, cap)` 用来指定长度和容量。`new(T)` 分配 `T` 并返回 `*T`，内存是零值。切片、map、channel 用 `make`，因为它们除了内存还要有运行时内部结构。

## map

```go
scores := map[string]int{"ada": 1}
v, ok := scores["missing"] // v 是 0，ok 是 false
```

用 `ok` 区分「没有这个键」和「值就是零值」。

map 的遍历顺序是随机的。不要依赖 `for range` 的顺序。并发读写同一个 map 会 panic。要并发就加锁，或者读 [Part C 的 sync 一章](../part-c/05-sync-atomic.md) 里 `sync.Map` 的适用条件。

Go 1.24 起，map 的实现换成 `internal/runtime/maps` 里基于 Swiss Table 的结构（组、控制字节、H1/H2）。这是实现。语言层面它仍是哈希表。

## 指针

`p := &x` 取地址，`*p` 解引用。Go 有指针，没有指针运算。零值指针是 `nil`，对 `nil` 解引用会 panic。

参数传递永远是把值复制进函数。复制一个指针只复制地址，所以函数可以通过指针改到调用方的变量。

## 运行

```bash
go run ./docs/go-tutorial/examples/parta/types
```

在 Go 1.27.1 上这次运行的输出是：

```text
rune=A code=U+0041
bytes=3 runes=1
slice=[1 2 3 4] len=4 cap=6
ada=1 ok=true
missing=false
ptr=11
```

`cap=6` 来自当前运行时的扩容（小切片按大约 2 倍增长，之后还会按大小类取整），**不是语言规范规定的常数**。换版本、换元素大小，容量都可以不同。公式在 [内存一章](../part-c/01-memory-escape.md) 里对照源码说明。

## 练习

1. 预测 `len("Go语言")` 和 `len([]rune("Go语言"))`，再写三行程序验证。
2. 把 `nums[:2]` 存下来，再对原切片 `append`，观察短切片的元素会不会变。容量够和不够各做一次。
