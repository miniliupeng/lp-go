# 05 可观测性

## 日志

结构化日志用 `log/slog`。JSON Handler 一行一个对象，字段是键值，不是拼进句子里的空格。

[obs 测试](../examples/partd/obs/slog_test.go) 把级别设成 Info：`Info("ready", "port", 8080)` 的输出里有 `"msg":"ready"` 和 `"port":8080`；`Debug` 不出现。

请求号放在 context 里，打日志时 `logger.With("request_id", id)`。见 [logctx](../examples/partb/logctx/logctx.go)。不要把密码、完整令牌、身份证号放进日志。

## 指标

`expvar` 适合进程内看一个数。要按时间聚合、按标签拆分，用 Prometheus 一类的系统。标准库没有 Prometheus 客户端。面试里讲清三类数就够：

- 速率：每秒请求数
- 错误：非 2xx 或业务失败的比例
- 耗时：延迟直方图，而不是只报一个平均值

Go 运行时自己的数在 `runtime/metrics`。不要用会 STW 的 `ReadMemStats` 做高频采集。内存上限的公式也在那一组指标里，见 [GC 一章](../part-c/03-gc.md)。

## 剖析端点

`net/http/pprof` 注册在默认 ServeMux 上。和生产 API 分开听，只绑 localhost，或者套一层认证。端点名称和 `goroutineleak` 的限制见 [性能一章](../part-c/07-pprof-trace-bench.md)。

追踪（trace）和日志不是一回事。日志是你写下来的事件；`runtime/trace` 是调度器的时间线。跨服务的 trace id 可以和请求号是同一个值，放在 context 里往下传。

## 运行

```bash
go test ./docs/go-tutorial/examples/partd/obs ./docs/go-tutorial/examples/partb/logctx -count=1
```

## 面试常问

- 为什么不用平均延迟：一次长尾会被大量快请求稀释。要分位。
- pprof 为什么不能暴露在公网上：它能看到栈、堆和运行中的程序状态。
- `slog` 的级别过滤发生在 Handler，调用处仍然会构造参数。极热的 Debug 要先看 `Enabled`，避免无用的字符串分配。

## 练习

把测试里的级别改成 Debug，确认 `"hidden"` 会出现。改完再改回去，保持仓库测试是 Info。
