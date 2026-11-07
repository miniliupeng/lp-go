# D1 HTTP 超时与关闭

路由写法见 [HTTP 服务](../能写服务/B3-HTTP服务.md)。这一章只讲 `http.Server` 的超时和 `Shutdown`。字段注释来自 Go 1.27.1 的 `net/http/server.go`。示例在 [examples/d1http](../../../examples/d1http/httpx_test.go)。

## 四个时间字段

零值或负数的含义：

| 字段 | 文档 |
| --- | --- |
| `ReadTimeout` | 读完整份请求，包含 body。0 或负数表示不设超时 |
| `ReadHeaderTimeout` | 只限制读请求头。0 时回退到 `ReadTimeout`。两者都是 0 或负数则没有超时 |
| `WriteTimeout` | 写响应的上限。读到下一个请求头时会重置。0 或负数表示不设超时 |
| `IdleTimeout` | keep-alive 等待下一个请求。0 时回退到 `ReadTimeout` |

测试构造的 `Server` 四个字段都是 0，也就是还没设置超时。大文件上传不要只靠一个很短的 `ReadTimeout`，它会把合法的 body 一起掐掉。头用 `ReadHeaderTimeout`，body 再由 Handler 自己决定。

`GET /files/{path...}` 会吃掉剩余路径。请求 `/files/a/b` 时 `PathValue("path")` 是 `a/b`。

## Shutdown

`Shutdown` 先关掉所有监听，再关掉空闲连接，然后等待还在处理的连接回到空闲。它不打断正在运行的 Handler。context 先到期时，返回的是 context 的错误。测试里 Handler 堵在 channel 上，100 毫秒的 context 得到 `context.DeadlineExceeded`。放开 Handler 之后，`Serve` 返回 `http.ErrServerClosed`。

文档还写明：调用 `Shutdown` 时，`Serve` 会立刻返回 `ErrServerClosed`。程序要等 `Shutdown` 返回，不要在它还没结束时让进程退出。被 Hijack 的连接（例如 WebSocket）不在等待范围内，要另外通知。`Shutdown` 之后这台 `Server` 不能再 `Serve`。

```bash
go test ./examples/d1http -count=1
```

## 练习

给测试里的 `Server` 设置 `ReadHeaderTimeout` 为 1 秒，确认它不再是零值。不要把这个改动留成一个会拖慢测试的真超时。
