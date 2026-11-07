# 07 pprof、trace 与基准

源码依据：Go 1.27.1 `runtime/pprof` 包注释里列出的预定义 profile。CPU profile 不在这张表里，它有单独的 API。

## 预定义 profile

| 名字 | 内容 |
| --- | --- |
| `goroutine` | 当前所有 goroutine 的栈 |
| `goroutineleak` | 被判定为泄漏的 goroutine 的栈。Go 1.27 起一般可用，不再是实验 |
| `allocs` | 历史上全部分配的采样 |
| `heap` | 仍存活对象的分配采样 |
| `threadcreate` | 创建 OS 线程时的栈 |
| `block` | 阻塞在同步原语上的栈 |
| `mutex` | 有竞争的 Mutex 持有者的栈 |

这些 profile 自己维护，对它们调用 `Add` 会 panic。`pprof.Lookup("goroutineleak")` 在本仓库的测试里必须非 nil。

泄漏的判定（Go 1.27 发布说明）：G 堵在 channel、Mutex、Cond 等同步原语上，而且从任何可运行的 goroutine 出发到不了这个原语，那么它不可能再被叫醒。GC 的可达性用来做这个判断。经过全局变量、或者仍可运行的 G 的局部变量还能摸到的阻塞，检测不到。HTTP 端点是 `/debug/pprof/goroutineleak`。`debug=2` 时格式接近未恢复的 panic。

导入 `net/http/pprof` 会在默认 `ServeMux` 上注册这些端点。不要把它暴露到公网。线上用单独的 localhost 端口。

profile 是采样。包注释写着：结果好不好取决于内核支持，并指向已知问题。不要把某一行的百分比读成精确的指令占比。

## CPU profile 与 trace

`pprof.StartCPUProfile` 写一块分析数据。测试把 50ms 的忙循环写进临时文件，并检查文件非空。这只证明「采样文件被写出来」，不证明某行代码的百分比。

`runtime/trace` 和 `go tool trace` 看的是调度、GC、系统调用的时间线，不是函数内部的指令热点。Go 1.27 起，`go tool trace -http=:6060` 只听 localhost；要听所有地址必须写成 `-http=0.0.0.0:6060`。

## 基准

```bash
go test -bench=BenchmarkAdd -benchmem -run=^$ ./docs/go-tutorial/examples/partc/profile
```

`b.N` 由测试框架调整。`-benchmem` 给出每次操作的分配字节和次数。把结果赋给包级变量，避免编译器把循环删掉。

**不要背某台机器上的 ns/op。** 换 CPU、换 Go 版本，数字都会变。比较时固定机器、固定版本，并看分配次数有没有从非零变成零。

基准里不要把一次性的初始化放进循环。需要计时时用 `b.ResetTimer()`。

## 排障顺序

内存涨了：先看是 `heap` 的 inuse 在涨，还是线程、cgo 或 `mmap` 不在 Go 堆里。Go 堆用 heap profile。`GOMEMLIMIT` 只约束运行时管理的内存，见 [GC 一章](03-gc.md)。

goroutine 数涨了：先看 `goroutine` 栈停在哪个 channel 或锁上，再用 `goroutineleak` 看有没有「永远叫不醒」的那一类。profile 为空不代表没有泄漏。

CPU 高：CPU profile 找热点，trace 看是不是 GC 或调度。两者回答的问题不同。

## 运行

```bash
go test ./docs/go-tutorial/examples/partc/profile -count=1
```

## 面试常问

- `allocs` 和 `heap` 的差别：一个是历史分配采样，一个是还活着的对象。
- `goroutineleak` 发现不了什么：阻塞点仍从全局变量或可运行 G 可达的情况。
- 为什么不能把 benchmark 数字写进设计文档当常数。

## 练习

给 `BenchmarkAdd` 加一个会 `append` 的对照基准，看 `-benchmem` 里 `allocs/op` 的差别。只比较同一轮命令的输出。
