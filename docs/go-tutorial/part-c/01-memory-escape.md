# 01 内存、逃逸与内存模型

版本：语言规则以 [Go memory model](https://go.dev/ref/mem)（2022-06-06 修订）为准。逃逸提示以 **Go 1.27.1** 编译器实际输出为准。

## 核心原理

Go 的局部变量可以在栈上，也可以在堆上。这不是由你写不写 `new` 决定的，而是编译器的逃逸分析：如果一个变量在函数返回后还可能被用到，它就不能只活在当前栈帧里。

返回地址是最直观的一种：

```go
func Leak(n int) *int { return &n }
```

Go 1.27.1 对它打印 `moved to heap: n`。更早的编译器常写 `escapes to heap`。两种措辞都表示：这个变量不能只放在栈上。测试同时接受这两种短语，避免把某一版的英文绑死。

`Stay` 只返回一个整数，局部量不需要逃出函数。它可以被内联。内联会改变逃逸结论，所以「我记得这个一定在栈上」不可靠，以 `-m` 输出为准。

切片扩容在 `runtime/slice.go` 的 `nextslicecap`：需要的新长度已经超过两倍时，直接用新长度；旧容量小于 256 时翻倍；否则每次增加 `(newcap + 3*256) >> 2`，从大约 2 倍平滑过渡到大约 1.25 倍。随后 `roundupsize` 还会按内存大小类取整，所以你在 `cap` 里看到的数可以比这个公式更大。[类型一章](../part-a/03-types.md) 里 `append` 一个 `int` 后 `cap=6`，就是这次运行的结果，不是规范常数。

## 内存模型

内存模型回答的是：一个 goroutine 里的读，什么时候一定能看见另一个 goroutine 里的写。

没有数据竞争的程序，表现得像所有 goroutine 被轮流放到一个处理器上执行。这叫 DRF-SC。数据竞争的定义是：对同一位置的并发访问里至少有一次写，而且这些访问没有 happens-before 关系，并且不是 `sync/atomic` 提供的原子访问。

happens-before 是「本 goroutine 里的顺序」加上「同步操作建立的 synchronized-before」的传递闭包。下面几条是规范里的原意，不是经验口诀：

- 启动 goroutine 的 `go` 语句，synchronized before 该 goroutine 开始执行。
- goroutine **退出**并不自动 synchronized before 程序里的任何其他事件。写完一个全局变量就让 goroutine 结束，别的 goroutine 不保证看得见。
- 对某个 channel 的发送，synchronized before 对应接收的完成。
- 关闭 channel，synchronized before 一次「因为关闭而收到零值」的接收。
- 从无缓冲 channel 接收，synchronized before 对应发送的完成。因此无缓冲 channel 可以用来做双向握手。缓冲容量为 1 时，这个保证消失。
- 容量为 C 的 channel 上，第 k 次接收 synchronized before 第 k+C 次发送的完成。
- 对同一把 `Mutex` 或 `RWMutex`，第 n 次 `Unlock`（n < m）synchronized before 第 m 次 `Lock` 返回。
- `once.Do(f)` 里 `f` 的返回，synchronized before 之后每一次 `Do` 的返回。
- 若原子操作 B 观察到原子操作 A 的效果，则 A synchronized before B。程序里全部原子操作看起来有一个顺序一致的总顺序。这和 C++ 的 sequentially consistent 原子、Java 的 `volatile` 是同一类语义。

有数据竞争时，Go 不像 C 那样允许编译器做任何事。实现可以报告竞争并终止；否则，对不超过一个机器字的读，必须观察到某个确实发生过、且还没被更晚的写覆盖的值。竞争仍然是错误。`-race` 用来在测试里把它们报出来。

## 实现要点

逃逸分析在编译期完成，运行时没有「再把堆上的对象搬回栈」这一步。接口、channel 发送、闭包、切片底层数组变大或被共享，都可能让变量逃逸。具体以编译器输出为准。

## 面试常问

- 逃逸是谁决定的：编译器，不是 `new` 这个关键字。
- 无缓冲 channel 为什么能传递「发送方的写已经完成」：接收 synchronized before 发送完成。
- goroutine 跑完了，main 为什么不一定看得见它的写入：退出本身不是同步事件。

## 对比陷阱

- 把「逃到堆上」说成「一定更慢所以不能用指针」。分配有成本，但是否逃逸要看 `-m`，不要凭感觉禁止指针。
- 背某个 `cap` 倍数。大小类取整之后，观察到的容量会变。
- 认为数据竞争在 Go 里是未定义行为、编译器可以删掉整个程序。规范限制了字长以内读的结果，但竞争程序仍然是错的。

## 可运行示例

```bash
go test ./docs/go-tutorial/examples/partc/escape -count=1
go build -gcflags=-m -o /dev/null ./docs/go-tutorial/examples/partc/escape
```

Go 1.27.1 上能看到 `can inline Leak` 和 `moved to heap: n`。

## 练习

给 `Leak` 再写一个「把 `n` 放进 `any` 再返回」的函数，重新跑 `-m`，看编译器怎么描述。
