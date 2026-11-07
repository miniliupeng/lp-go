# 01 net/http

源码依据：Go 1.27.1 `net/http/server.go` 里 `Server` 各超时字段的注释，以及 `Shutdown` 的注释。路由行为在 Go 1.22 及以后。

## 超时字段

零值的含义以字段注释为准，测试 [route_test.go](../examples/partd/httpx/route_test.go) 只锁定「新 Server 的这四个字段都是 0」。

| 字段 | 为 0 或负数时 | 正数时 |
| --- | --- | --- |
| `ReadTimeout` | 不限制读完整请求（含 body）的时间 | 从接受连接到读完 body 的上限 |
| `ReadHeaderTimeout` | 为 0 时改用 `ReadTimeout`。两者都不是正数则没有读头超时 | 只限制读请求头。读完头会重置连接的读截止时间，Handler 可以自己决定 body 有多慢 |
| `WriteTimeout` | 不限制写响应 | 从读完请求头开始算的写上限。读到下一个请求头时会重置。它不让 Handler 按请求单独放宽 |
| `IdleTimeout` | 为 0 时改用 `ReadTimeout`。两者都不是正数则没有空闲超时 | keep-alive 下等待下一个请求的时间 |

文档建议多数服务设置 `ReadHeaderTimeout`，因为 `ReadTimeout` 把 body 也算进去，Handler 没法按请求放宽。四个字段可以同时设。

`ReadTimeout` 包住整个 body 时，慢上传会把连接掐掉。大上传应把 `ReadTimeout` 留空或设得很大，用 `ReadHeaderTimeout` 防慢速头，再用 `http.MaxBytesReader` 限制 body 大小。

## 路由

`http.NewServeMux` 的模式可以带方法：`GET /files/{path...}`。`r.PathValue("path")` 得到剩余路径。测试请求 `/files/a/b`，正文是 `a/b`。

错误的方法返回 405，见 [HTTP 服务一章](../part-b/02-http-service.md)。

## Shutdown

`Server.Shutdown(ctx)` 的注释：

- 先关闭所有 Listener，再关闭空闲连接，然后等待仍在处理的连接回到空闲再关掉。
- 它不打断正在执行的 Handler。Handler 不返回，`Shutdown` 就等到 `ctx` 结束并返回 `ctx.Err()`。
- `Serve` / `ListenAndServe` 在 `Shutdown` 被调用后返回 `http.ErrServerClosed`。进程要等 `Shutdown` 返回，不要在 `Serve` 一返回就 `os.Exit`。
- 被 Hijack 的连接（例如 WebSocket）既不会被关掉，也不会被等待。要另外通知它们。`RegisterOnShutdown` 用来挂收尾钩子。

[shutdown_test.go](../examples/partd/shutdown/shutdown_test.go) 有两个测试：没有在途请求时 `Shutdown` 成功，`Serve` 得到 `ErrServerClosed`；Handler 堵住时，100ms 的 context 让 `Shutdown` 返回 `context.DeadlineExceeded`。第二个测试在断言之后才放开 Handler，避免把「Handler 已经跑完」误判成超时。

## 客户端连接复用

`Transport` 按 host 保留空闲连接。`MaxIdleConns` 是总数，`MaxIdleConnsPerHost` 默认很小（2）。高并发客户端不调大会把连接不断建拆。

Go 1.27 的发布说明：HTTP/1 的 `Response.Body` 在 `Close` 时会把未读完的内容排掉，有一个保守上限，以便复用连接。大多数程序这是加速。如果本来就靠「不读 body、让连接堆着」这种错误用法，再把 `MaxIdleConns` 设成 0 或每个请求换一个 `Client`，可能变慢。那种情况应设置 `DisableKeepAlives`，并回头检查是不是客户端用错了。

HTTP/2 服务端在 Go 1.27 接受 RFC 9218 的客户端优先级。要恢复按流轮询，把 `Server.DisableClientPriority` 设为 true。

## 面试常问

- 四个超时各自从什么时候开始算。
- `Shutdown` 会不会取消请求的 context、会不会等 WebSocket：按上面的文档回答，不要说它会把每个 Handler 的 context 取消掉。文档说的是不打断活跃连接。
- 为什么必须 `Close` 响应体：不关会占着连接，连接池里的空闲连接补不上。

## 运行

```bash
go test ./docs/go-tutorial/examples/partd/httpx ./docs/go-tutorial/examples/partd/shutdown -count=1
```

## 练习

给一个测试用的 `Server` 只设置 `ReadHeaderTimeout`，读源码注释说明 `IdleTimeout` 此时会落到哪个值。
