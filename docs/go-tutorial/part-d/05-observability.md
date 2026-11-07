# 05 可观测性

## 这一章你要会什么

分清日志、指标和剖析各回答什么问题，并用本仓库的测试核对 `slog` 与 `runtime/metrics` 的行为。

## 日志

结构化日志用 `log/slog`。JSON Handler 一行一个对象。字段是键值，不要把 `port=8080` 拼进消息句子里。

[obs 测试](../examples/partd/obs/slog_test.go) 把 Handler 的级别设成 Info：

- `Info("ready", "port", 8080)` 的输出里有 `"msg":"ready"` 和 `"port":8080`。
- `Debug("hidden")` 不会出现在缓冲区里。
- `logger.Enabled(ctx, slog.LevelDebug)` 为 false，`LevelInfo` 为 true。

级别过滤发生在 Handler。调用 `Debug` 时，参数仍然会先被求值。极热路径上如果拼一条 Debug 消息很贵，先看 `Enabled`，为 false 就不要做那次分配。测试锁住的是 `Enabled` 的返回值；它没有证明「没调用 Enabled 时参数不会求值」，那是 Go 的函数调用规则，参数在进入函数前就算好了。

请求号放在 context 里，打日志时 `logger.With("request_id", id)`。见 [logctx](../examples/partb/logctx/logctx.go)。不要把密码、完整令牌、身份证号放进日志。日志是你写下来的事件，不是调度器的时间线。

## 指标

面试里把对外的服务先收成三类数：

- **速率**：每秒请求数。
- **错误**：非 2xx，或者你定义的业务失败，占请求的比例。
- **耗时**：延迟的直方图。平均值会被大量快请求稀释，一次长尾在平均数里几乎看不见。要的是分位，例如 p99。

标准库里的 `expvar` 适合在进程内看一个当前值。要按时间聚合、按标签拆分，用 Prometheus 一类的系统。标准库没有 Prometheus 客户端。标签的取值必须是有限集合：把用户 id、URL 全文放进标签，时间序列的数量会跟着用户数涨。

Go 运行时自己的数在 `runtime/metrics`。`metrics.Read` 按 `Sample.Name` 填充 `Value`。名字不在 `metrics.All()` 里时，`Value.Kind()` 是 `KindBad`，这时不能调用 `Uint64()`，那个方法在种类不对时会 panic。测试读取两个已知名字：

- `/memory/classes/total:bytes`：正常进程里是非零的 `uint64`。
- `/memory/classes/heap/released:bytes`：种类必须有效，数值可以是 0。

两者相减就是 [GC 一章](../part-c/03-gc.md) 里内存软上限使用的那组量，对应 `MemStats.Sys - MemStats.HeapReleased`。同一块 `Sample` 切片可以交给下一次 `Read`。一次 `Read` 还没返回时，不要并发读里面的 `Value`。`runtime.ReadMemStats` 会短暂 STW，不要用它做高频采集。

## 剖析和追踪

`net/http/pprof` 把剖析端点注册在默认的 `ServeMux` 上。它能看到栈、堆和运行中的程序状态，所以和生产 API 分开听，只绑 localhost，或者套一层认证。端点名称和 `goroutineleak` 能发现什么，见 [性能一章](../part-c/07-pprof-trace-bench.md)。

`runtime/trace` 记录的是调度、GC 和系统调用的时间线，用来解释「这 100 毫秒卡在哪」。跨服务传递的 trace id 可以和请求号是同一个值，放在 context 里往下传。日志里的请求号和 trace 里的 id 对上，才能从一条错误日志跳到那次请求的时间线。

## 运行

```bash
go test ./docs/go-tutorial/examples/partd/obs ./docs/go-tutorial/examples/partb/logctx -count=1
```

## 面试常问

- 为什么延迟看分位：平均数会被大量快请求拉低。
- pprof 为什么不能暴露在公网上：调用者能看到栈和堆。
- 热路径上的 Debug 为什么先看 `Enabled`：`Debug` 的参数在调用前就会求值。
- 内存软上限看哪两个 `runtime/metrics` 名字，以及为什么不用 `ReadMemStats` 打点。

## 练习

把测试里的 Handler 级别改成 `slog.LevelDebug`，确认 `"hidden"` 会出现，并且 `Enabled(..., slog.LevelDebug)` 变成 true。看完再改回 Info，保持仓库里的测试是现在的断言。
