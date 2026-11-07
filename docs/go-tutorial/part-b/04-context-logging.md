# 04 context 与日志

## context 做什么

`context.Context` 传递取消、截止时间和少量请求范围的数据。它不是可选参数包。

- `WithCancel`：调用 `cancel` 后，`ctx.Done()` 关闭，`ctx.Err()` 是 `context.Canceled`。
- `WithTimeout` / `WithDeadline`：到点后 `ctx.Err()` 是 `context.DeadlineExceeded`。`cancel` 仍然要调用，用来释放定时器资源。
- `WithCancelCause`（Go 1.20）：取消时可以附带原因，用 `context.Cause(ctx)` 取出。普通的 `ctx.Err()` 仍是 `Canceled`。

取消 context **不会**把正在跑的函数杀掉。函数要自己在循环里看 `ctx.Done()`，或者把 `ctx` 传给 `database/sql`、`net/http` 这类会监听它的 API。

不要把 `context.Background()` 或 `TODO()` 用在已经有请求 context 的处理函数里。入口（`main`、测试、传入的请求）才创建根 context。

`WithValue` 的键用未导出的空结构体类型，避免和别的包撞键。本教程只用它放请求号，见 [logctx](../examples/partb/logctx/logctx.go)。业务配置、数据库句柄放在结构体字段里。

## slog

`log/slog` 从 Go 1.21 起在标准库里。`slog.NewJSONHandler` 打出一行一个 JSON 对象。级别低于 `HandlerOptions.Level` 的记录会被丢掉。测试在 [obs](../examples/partd/obs/slog_test.go)：`Info` 保留，`Debug` 不出现。

给 logger 加字段用 `logger.With("request_id", id)`，不要把请求号拼进消息字符串里。

## 运行

```bash
go test ./docs/go-tutorial/examples/partb/logctx ./docs/go-tutorial/examples/partd/obs -count=1
```

## 面试常问

- 超时了 goroutine 还会不会停：只有它自己观察 context，或者它调用的库观察 context，它才会停。
- `Err()` 和 `Cause()` 的差别：前者是 `Canceled` 或 `DeadlineExceeded`，后者是你附加的原因。
- context 值适合放什么：请求号、追踪号这类请求范围的元数据。

## 练习

写一个循环，每步检查 `ctx.Done()`。用已经取消的 context 调用它，确认函数返回 `ctx.Err()`。
